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

func (r *ReserveRequest) Save(db *gorm.DB) {
	db.FirstOrCreate(&r)
}

func (r *ReserveRequest) Update(db *gorm.DB) {

}

func (r *ReserveRequest) Delete(db *gorm.DB) {
	db.Delete(r)
}
