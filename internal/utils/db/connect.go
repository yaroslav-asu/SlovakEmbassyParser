package db

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"main/internal/utils/vars"
	gorm_models "main/models/gorm"
)

func Connect() *gorm.DB {
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", vars.DbUser, vars.DbPassword, vars.DbName)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		zap.L().Error("Failed to connect db")
		zap.L().Info("Trying to reconnect db")
		return Connect()
	}
	err = db.AutoMigrate(&gorm_models.Reservation{}, &gorm_models.City{}, &gorm_models.ReserveRequest{}, &gorm_models.User{})
	if err != nil {
		zap.L().Error("failed to auto migrate database")
		zap.L().Warn("Continuing without auto migration")
	}
	return db
}
