package session

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"main/internal/session/captcha"
	"main/internal/session/captcha/rucaptcha"
	"main/internal/utils/funcs"
	"main/internal/utils/vars"
	"net/url"
	"os"
	"strconv"
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
	s.SendCaptchaToSolve()
	s.captcha.DeleteCaptcha()
	time.Sleep(CaptchaSolveWaitTime)
	solve := s.GetCaptchaSolve()
	zap.L().Info("Got captcha solve: " + solve)
	return solve
}

func (s *Session) SolveNewCaptcha() string {
	s.DownloadCaptcha()
	return s.solveCaptcha()
}
