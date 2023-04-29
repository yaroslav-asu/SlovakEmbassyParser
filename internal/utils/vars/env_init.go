package vars

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

var (
	SiteUrl                 string
	DefaultUserName         string
	DefaultUserPassword     string
	RunningMode             string
	RuCaptchaApiKey         string
	DbUser                  string
	DbPassword              string
	DbName                  string
	CaptchaSolveProjectPath string
)

func InitDefaultEnv() {
	zap.L().Info("Started to initialize environmental vars")
	err := godotenv.Load(".env")
	if err != nil {
		zap.L().Fatal("Failed to load .env file")
	}
	SiteUrl = os.Getenv("SITE_URL")
	DefaultUserName = os.Getenv("DEFAULT_USER_NAME")
	DefaultUserPassword = os.Getenv("DEFAULT_USER_PASSWORD")
	RunningMode = os.Getenv("RUNNING_MODE")
	RuCaptchaApiKey = os.Getenv("RUCAPTCHA_API_KEY")
	CaptchaSolveProjectPath = os.Getenv("CAPTCHA_SOLVE_PROJECT_PATH")
	zap.L().Info("Finished initializing environmental vars")
}

func InitDbEnv() {
	zap.L().Info("Started to initialize environmental vars for db")
	err := godotenv.Load(".env.db")
	if err != nil {
		zap.L().Fatal("Failed to load .env.db file")
	}
	DbUser = os.Getenv("POSTGRES_USER")
	DbPassword = os.Getenv("POSTGRES_PASSWORD")
	DbName = os.Getenv("POSTGRES_DB")
	zap.L().Info("Finished initializing environmental vars for db")
}

func InitEnv() {
	InitDefaultEnv()
	InitDbEnv()
}
