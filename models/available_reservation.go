package models

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AvailableReservation struct {
	CityId string
	Date   Date
}

func (a AvailableReservation) SaveToDB(db *gorm.DB) {
	zap.L().Info("Saved to DB")
	db.FirstOrCreate(&a, a)
}

func (a AvailableReservation) DeleteFromDB(db *gorm.DB) {
	db.Where("date = ? and id = ?", a.Date, a.CityId).Delete(&a)
}

type AvailableReservations []DbModel

func (a AvailableReservations) SaveToDB(db *gorm.DB) {
	for reservationId := range a {
		a[reservationId].SaveToDB(db)
	}
}
