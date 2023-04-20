package captcha

import (
	"encoding/base64"
	"fmt"
	"go.uber.org/zap"
	"os"
)

type Captcha struct {
	title       string
	RucaptchaId string
}

func NewCaptcha(title string) Captcha {
	return Captcha{
		title: title,
	}
}
func (c Captcha) Path() string {
	return fmt.Sprintf("captcha/%s.png", c.title)
}

func (c Captcha) Base64() string {
	bytes, err := os.ReadFile(c.Path())
	if err != nil {
		zap.L().Error("Failed to open ")
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

func (c Captcha) SolveCaptchaOffline() string {
	var textCaptcha string
	fmt.Print(fmt.Sprintf("Type captcha solve %s: ", c.title))
	_, err := fmt.Scan(&textCaptcha)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return textCaptcha
}

func (c Captcha) DeleteCaptcha() {
	err := os.Remove(c.Path())
	if err != nil {
		zap.L().Warn("Failed to delete captcha: '" + c.title + "'")
	}
}
