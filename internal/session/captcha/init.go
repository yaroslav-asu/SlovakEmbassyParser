package captcha

import (
	"go.uber.org/zap"
	"os"
)

func Init() {
	dirs := []string{"captcha", "captcha/right", "captcha/wrong"}
	for _, dirName := range dirs {
		if _, err := os.Stat(dirName); err != nil {
			if err := os.Mkdir(dirName, os.ModePerm); err != nil {
				zap.L().Fatal("Failed to create folder: " + dirName)
			}
		}
	}
}
