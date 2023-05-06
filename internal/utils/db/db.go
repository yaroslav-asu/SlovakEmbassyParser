package db

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main/internal/utils/vars"
	gorm_models "main/models/gorm"
	"time"
)

var reconnectTime = 5 * time.Second

func Init() {
	db := Connect()
	defer Close(db)
	err := db.AutoMigrate(
		&gorm_models.Reservation{},
		&gorm_models.City{},
		&gorm_models.ReserveRequest{},
		&gorm_models.User{},
		&gorm_models.Proxy{},
		&gorm_models.City{},
		&gorm_models.DayCell{},
		&gorm_models.Cookie{},
		&gorm_models.Month{},
	)
	if err != nil {
		zap.L().Error("failed to auto migrate database")
		zap.L().Info("Continuing without auto migration")
	}
}

func Connect() *gorm.DB {
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", vars.DbUser, vars.DbPassword, vars.DbName)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		zap.L().Error("Failed to connect db")
		zap.L().Info("Trying to reconnect db")
		time.Sleep(reconnectTime)
		reconnectTime *= 2
		return Connect()
	}
	return db
}

func Close(db *gorm.DB) {
	postgresDB, err := db.DB()
	if err != nil {
		zap.L().Error("Failed to get db instance: " + err.Error())
		zap.L().Info("DB connection wasn't close")
		return
	}
	err = postgresDB.Close()
	if err != nil {
		zap.L().Info("DB connection wasn't close: " + err.Error())
	}
}
