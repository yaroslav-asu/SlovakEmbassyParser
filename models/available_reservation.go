package models

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AvailableReservation struct {
	gorm.Model
	Date   Date
	CityId string
}

func (a AvailableReservation) SaveToDb(db *gorm.DB) {
	zap.L().Info("Saved to DB")
	db.FirstOrCreate(&a, a)
}

type AvailableReservations []DbModel
