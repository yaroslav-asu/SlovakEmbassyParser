package session

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"main/internal/session/captcha"
	"main/internal/session/captcha/rucaptcha"
	"main/internal/utils/db"
	"main/internal/utils/funcs"
	"main/internal/utils/vars"
	gorm_models "main/models/gorm"
	"main/models/gorm/datetime"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const CaptchaSolveWaitTime = 2 * time.Second

func (s *Session) DownloadCaptcha() captcha.Captcha {
	zap.L().Info("Starting to download captcha of user: " + s.User.UserName)
	res := s.Get(funcs.Linkify("simpleCaptcha.png"))
	s.captcha = captcha.NewCaptcha(s.User.UserName)
	file, err := os.Create(s.captcha.Path())
	if err != nil {
		zap.L().Error("Can't create captcha.png with error: " + err.Error())
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		zap.L().Error("Can't write captcha bytes to file with error: " + err.Error())
	}
	return s.captcha
}

func (s *Session) SendCaptchaToSolve() {
	zap.L().Info("Sending post request with captcha: " + s.captcha.Format())
	res := s.PostForm(
		"http://rucaptcha.com/in.php",
		url.Values{
			"key":      {vars.RuCaptchaApiKey},
			"method":   {"base64"},
			"body":     {s.captcha.Base64()},
			"phrase":   {"0"},
			"regsense": {"1"},
			"language": {"2"},
			"json":     {"1"},
		},
	)
	response := rucaptcha.ParseRucaptchaResponse(res)
	if response.Status != 1 {
		zap.L().Error("Captcha didn't send with error title: " + response.Request + " and error text: " + response.ErrorText)
	}
	zap.L().Info("Setting captcha id of rucaptcha server to session")
	s.captcha.RucaptchaId = response.Request
	defer res.Body.Close()
}

func (s *Session) GetCaptchaSolve() string {
	for ; ; time.Sleep(CaptchaSolveWaitTime) {
		zap.L().Info("Trying to get solved captcha: " + s.captcha.Format())
		res := s.Get(fmt.Sprintf("http://rucaptcha.com/res.php?key=%s&action=get&id=%s&json=1", vars.RuCaptchaApiKey, s.captcha.RucaptchaId))
		response := rucaptcha.ParseRucaptchaResponse(res)
		switch response.Status {
		case 0:
			zap.L().Info("Captcha isn't ready, next try will be after " + strconv.Itoa(int(CaptchaSolveWaitTime/time.Second)) + " seconds")
		case 1:
			return response.Request
		default:
			zap.L().Error("Got unknown response status: " + response.Format())
			zap.L().Info("Breaking checking response cycle")
			return ""
		}
	}
}

func (s *Session) solveCaptcha() string {
	solve, err := s.captcha.PredictSolve()
	if err != nil {
		zap.L().Info("Error in prediction captcha solve, starting to get solve from rucaptcha")
		s.SendCaptchaToSolve()
		time.Sleep(CaptchaSolveWaitTime)
		solve = s.GetCaptchaSolve()
	}
	s.captcha.Rename(solve)
	zap.L().Info("Got captcha solve: " + solve)
	return solve
}

func (s *Session) SolveNewCaptcha() string {
	unproxiedSession := NewBlankSession()
	unproxiedSession.captcha = s.DownloadCaptcha()
	solve := unproxiedSession.solveCaptcha()
	s.captcha.Solve = solve
	return s.captcha.Solve
}

func (s *Session) CheckCaptchaSolve() bool {
	now := datetime.NewDateYMD(2023, 5, 10)
	consularId := "590"
	res := s.PostForm(
		funcs.Linkify(fmt.Sprintf("calendarDay.do?day=%s&calendarId=&consularPostId=%s", now.Format(datetime.DateOnly), consularId)),
		url.Values{
			"calendar.timeOfVisit":               {now.Format(datetime.DateOnly)},
			"calendar.sequenceNo":                {""},
			"calendar.timeSlot.timeSlotId":       {""},
			"calendar.consularPost.consularPost": {""},
			"calendar.calendarId":                {""},
			"captcha":                            {s.captcha.Solve},
		},
	)
	root, err := funcs.ResponseToSoup(res)
	if err != nil {
		zap.L().Info("Failed to read check captcha response, quitting with false value")
		return false
	}
	for _, el := range root.FindAll("script") {
		if strings.Contains(el.FullText(), "captcha") {
			errParts := strings.Split(funcs.StripString(el.FullText()), ",")
			zap.L().Info("Captcha solve isn't right: " + errParts[len(errParts)-1])
			return false
		}
	}
	zap.L().Info("Captcha solve is right")
	return true
}

func (s *Session) ParseCaptchas(count int) {
	for i := 0; i < count; i++ {
		solve := s.SolveNewCaptcha()
		var dirName string
		if s.CheckCaptchaSolve() {
			dirName = "right"
		} else {
			dirName = "wrong"
		}
		err := os.Rename(fmt.Sprintf("captcha/%s.png", solve), fmt.Sprintf("captcha/%s/%s.png", dirName, solve))
		if err != nil {
			log.Fatal(err)
		}
		funcs.SleepTime(5, 20)
	}
}

func StartParseCaptchas() {
	d := db.Connect()
	defer db.Close(d)
	var u gorm_models.User
	d.Model(&gorm_models.User{}).Where("id = 2").First(&u)
	s := NewLoggedInSession(u.UserName, u.Password)
	s.ParseCaptchas(100)
}
