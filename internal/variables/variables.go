package variables

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

var SiteUrl string
var DefaultUserName string
var DefaultUserPassword string

func InitEnv() {
	zap.L().Info("Started to initialize environmental variables")
	err := godotenv.Load(".env")
	if err != nil {
		zap.L().Warn("Failed to load .env file")
	}
	SiteUrl = os.Getenv("SITE_URL")
	DefaultUserName = os.Getenv("DEFAULT_USER_NAME")
	DefaultUserPassword = os.Getenv("DEFAULT_USER_PASSWORD")
	zap.L().Info("Environmental variables successfully initialized")
}
