package parser

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"main/internal/utils/funcs"
	"os"
)

func (p *Parser) DownloadCaptcha() {
	res := p.get(funcs.Linkify("simpleCaptcha.png"))
	defer res.Body.Close()
	file, err := os.Create("captcha.png")
	if err != nil {
		zap.L().Error("Cant create captcha.png with error: " + err.Error())
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		zap.L().Error("Cant write captcha bytes to file with error: " + err.Error())
	}
}

func (p *Parser) SolveCaptcha() string {
	var textCaptcha string
	_, err := fmt.Scan(textCaptcha)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return textCaptcha
}
