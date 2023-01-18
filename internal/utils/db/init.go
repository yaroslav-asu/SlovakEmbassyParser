package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main/internal/utils/variables"
	"main/models"
)

func Connect() *gorm.DB {
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", variables.DbUser, variables.DbPassword, variables.DbName)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&models.City{})
	if err != nil {
		panic("failed to auto migrate database")
	}
	return db
}
