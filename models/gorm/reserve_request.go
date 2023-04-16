package gorm

import (
	"gorm.io/gorm"
	"main/models/gorm/datetime"
)

type ReserveRequest struct {
	Id     uint `gorm:"primaryKey"`
	UserId int
	User   User
	CityId string
	City   City
	Start  datetime.Date
	End    datetime.Date
}

func (r *ReserveRequest) SaveToDB(db *gorm.DB) {
	db.FirstOrCreate(&r)
}

func (r *ReserveRequest) DeleteFromDB(db *gorm.DB) {
	db.Delete(r)
}
