package gorm

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/models"
	"main/models/gorm/datetime"
)

type Reservation struct {
	Id       uint `gorm:"primaryKey"`
	CityId   string
	City     City
	DateTime datetime.Date
}

func (r Reservation) Save(db *gorm.DB) {
	zap.L().Info("Saved to DB")
	db.FirstOrCreate(&r, r)
}

func (r Reservation) Update(db *gorm.DB) {

}

func (r Reservation) Delete(db *gorm.DB) {
	db.Where("date = ? and id = ?", r.DateTime, r.CityId).Delete(&r)
}

type Reservations []models.DBModel

func (r Reservations) SaveToDB(db *gorm.DB) {
	for reservationId := range r {
		r[reservationId].Save(db)
	}
}

func (r Reservations) DeleteFromDB(db *gorm.DB) {
	for reservationId := range r {
		r[reservationId].Delete(db)
	}
}
