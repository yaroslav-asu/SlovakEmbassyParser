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
	"time"
)

const CaptchaSolveWaitTime = 2 * time.Second

func (s *Session) DownloadCaptcha() captcha.Captcha {
	res := s.Get(funcs.Linkify("simpleCaptcha.png"))
	newCaptcha := captcha.NewCaptcha(s.username)
	file, err := os.Create(newCaptcha.Path())
	if err != nil {
		zap.L().Error("Cant create captcha.png with error: " + err.Error())
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		zap.L().Error("Cant write captcha bytes to file with error: " + err.Error())
	}
	s.captcha = newCaptcha
	return newCaptcha
}

func (s *Session) SendCaptchaToSolve() {
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
		zap.L().Error("Image dont sent with error title: " + response.Request + " and error text: " + response.ErrorText)
	}
	s.captcha.RucaptchaId = response.Request
	defer res.Body.Close()
}

func (s *Session) GetCaptchaSolve() string {
	for ; ; time.Sleep(CaptchaSolveWaitTime) {
		res := s.Get(fmt.Sprintf("http://rucaptcha.com/res.php?key=%s&action=get&id=%s&json=1", vars.RuCaptchaApiKey, s.captcha.RucaptchaId))
		response := rucaptcha.ParseRucaptchaResponse(res)
		if response.Status == 1 {
			return response.Request
		}
	}
}

func (s *Session) solveCaptcha() string {
	s.SendCaptchaToSolve()
	s.captcha.DeleteCaptcha()
	time.Sleep(CaptchaSolveWaitTime)
	return s.GetCaptchaSolve()
}

func (s *Session) SolveNewCaptcha() string {
	s.DownloadCaptcha()
	return s.solveCaptcha()
}
