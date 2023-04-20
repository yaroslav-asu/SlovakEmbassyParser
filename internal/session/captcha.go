package session

import (
	"go.uber.org/zap"
	"io"
	"main/internal/session/captcha"
	"main/internal/utils/funcs"
	"os"
)

func (s *Session) DownloadCaptcha() captcha.Captcha {
	res := s.Get(funcs.Linkify("simpleCaptcha.png"))
	file, err := os.Create("captcha.png")
	if err != nil {
		zap.L().Error("Cant create captcha.png with error: " + err.Error())
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		zap.L().Error("Cant write captcha bytes to file with error: " + err.Error())
	}
	return captcha.NewCaptcha(s.username)
}
