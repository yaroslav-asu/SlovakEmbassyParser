package db

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main/internal/utils/vars"
	gorm_models "main/models/gorm"
)

func Connect() *gorm.DB {
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", vars.DbUser, vars.DbPassword, vars.DbName)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("failed to connect database")
	}
	err = db.AutoMigrate(&gorm_models.Reservation{}, &gorm_models.City{}, &gorm_models.ReserveRequest{}, &gorm_models.User{})
	if err != nil {
		zap.L().Fatal("failed to auto migrate database")
	}
	return db
}
