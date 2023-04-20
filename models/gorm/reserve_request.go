package gorm

import (
	"main/internal/utils/db"
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

func (r *ReserveRequest) SaveToDB(db *db.DB) {
	db.DB.FirstOrCreate(&r)
}

func (r *ReserveRequest) DeleteFromDB(db *db.DB) {
	db.DB.Delete(r)
}
