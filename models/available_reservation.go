package models

import (
	"gorm.io/gorm"
)

type AvailableReservation struct {
	gorm.Model
	Date   Date
	CityId string
}

func (a AvailableReservation) FirstOrCreate(db *gorm.DB) {
	db.FirstOrCreate(&a, a)
}

type AvailableReservations struct {
	Reservations []DbModel
}

func (a AvailableReservations) SaveToDb(db *gorm.DB) {
	for _, reservation := range a.Reservations {
		reservation.FirstOrCreate(db)
	}
}
