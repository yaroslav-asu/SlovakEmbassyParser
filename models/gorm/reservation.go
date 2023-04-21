package gorm

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/models"
	"main/models/gorm/datetime"
)

type Reservation struct {
	Id     uint `gorm:"primaryKey"`
	CityId string
	City   City
	Date   datetime.Date
}

func (r Reservation) SaveToDB(db *gorm.DB) {
	zap.L().Info("Saved to DB")
	db.FirstOrCreate(&r, r)
}

func (r Reservation) DeleteFromDB(db *gorm.DB) {
	db.Where("date = ? and id = ?", r.Date, r.CityId).Delete(&r)
}

type Reservations []models.DbModel

func (r Reservations) SaveToDB(db *gorm.DB) {
	for reservationId := range r {
		r[reservationId].SaveToDB(db)
	}
}

func (r Reservations) DeleteFromDB(db *gorm.DB) {
	for reservationId := range r {
		r[reservationId].DeleteFromDB(db)
	}
}
