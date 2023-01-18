package variables

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

var (
	SiteUrl             string
	DefaultUserName     string
	DefaultUserPassword string
	RunningMode         string
	DbUser              string
	DbPassword          string
	DbName              string
)

func InitDefaultEnv() {
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

func InitDbEnv() {
	zap.L().Info("Started to initialize environmental variables for db")
	err := godotenv.Load("db.env")
	if err != nil {
		zap.L().Warn("Failed to load db.env file")
	}
	DbUser = os.Getenv("POSTGRES_USER")
	DbPassword = os.Getenv("POSTGRES_PASSWORD")
	DbName = os.Getenv("POSTGRES_DB")
}

func InitEnv() {
	InitDefaultEnv()
	InitDbEnv()
}
