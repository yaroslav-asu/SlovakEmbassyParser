package captcha

import (
	"fmt"
	"go.uber.org/zap"
)

type Captcha struct {
	title string
}

func NewCaptcha(title string) Captcha {
	return Captcha{
		title: title,
	}
}

func (c Captcha) SolveCaptcha() string {
	var textCaptcha string
	fmt.Print(fmt.Sprintf("Type captcha solve %s: ", c.title))
	_, err := fmt.Scan(&textCaptcha)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return textCaptcha
}
