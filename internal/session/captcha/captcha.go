package captcha

import (
	"encoding/base64"
	"fmt"
	"go.uber.org/zap"
	"main/internal/utils/vars"
	"os"
	"os/exec"
	"path/filepath"
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

func (c Captcha) SolveByInput() string {
	var textCaptcha string
	fmt.Print(fmt.Sprintf("Type captcha solve %s: ", c.title))
	_, err := fmt.Scan(&textCaptcha)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return textCaptcha
}

func (c Captcha) PredictSolve() (string, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		zap.L().Error("Failed to get current working directory")
		return "", err
	}
	out, err := exec.Command(
		filepath.Join(vars.CaptchaSolveProjectPath, "venv/bin/python"),
		filepath.Join(vars.CaptchaSolveProjectPath, "solver.py"),
		filepath.Join(rootDir, fmt.Sprintf("captcha/%s.png", c.title)),
	).Output()
	if err != nil {
		zap.L().Error("Failed to get predict form model")
		return "", err
	}
	return string(out), nil
}

func (c Captcha) Format() string {
	return fmt.Sprintf("Captcha{title: %s, RucaptchaId: %s}", c.title, c.RucaptchaId)
}

func (c Captcha) Delete() {
	zap.L().Info("Starting to delete captcha: " + c.Format())
	err := os.Remove(c.Path())
	if err != nil {
		zap.L().Warn("Failed to delete captcha: " + c.Format())
	}
}

func (c Captcha) Rename(newName string) {
	zap.L().Info("Renaming captcha from: " + c.title + " to: " + newName)
	err := os.Rename(fmt.Sprintf("captcha/%s.png", c.title), fmt.Sprintf("captcha/%s.png", newName))
	if err != nil {
		zap.L().Error("Failed to rename captcha")
	}
}
