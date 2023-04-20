package gorm

import (
	"go.uber.org/zap"
	"main/internal/utils/db"
	"main/models"
	"main/models/gorm/datetime"
)

type Reservation struct {
	Id     uint `gorm:"primaryKey"`
	CityId string
	City   City
	Date   datetime.Date
}

func (r Reservation) SaveToDB(db *db.DB) {
	zap.L().Info("Saved to DB")
	db.DB.FirstOrCreate(&r, r)
}

func (r Reservation) DeleteFromDB(db *db.DB) {
	db.DB.Where("date = ? and id = ?", r.Date, r.CityId).Delete(&r)
}

type Reservations []models.DbModel

func (r Reservations) SaveToDB(db *db.DB) {
	for reservationId := range r {
		r[reservationId].SaveToDB(db)
	}
}

func (r Reservations) DeleteFromDB(db *db.DB) {
	for reservationId := range r {
		r[reservationId].DeleteFromDB(db)
	}
}
