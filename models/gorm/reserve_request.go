package gorm

import (
	"gorm.io/gorm"
	"main/internal/datetime"
)

type ReserveRequest struct {
	Id     uint `gorm:"primaryKey"`
	UserId int
	CityId string
	Start  datetime.Date
	End    datetime.Date
}

func (r ReserveRequest) SaveToDB(db *gorm.DB) {
	db.FirstOrCreate(&r)
}

func (r ReserveRequest) DeleteFromDB(db *gorm.DB) {

}
