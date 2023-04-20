package captcha

import (
	"go.uber.org/zap"
	"os"
)

func InitCaptcha() {
	if _, err := os.Stat("captcha"); err != nil {
		if err := os.Mkdir("captcha", os.ModePerm); err != nil {
			zap.L().Fatal("Failed to create folder captcha")
		}
	}
}
