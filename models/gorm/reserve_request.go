package gorm

import (
	"gorm.io/gorm"
)

type ReserveRequest struct {
	Id      uint `gorm:"primaryKey"`
	UserId  int
	User    User
	CityId  string
	City    City
	MonthId uint
	Month   Month
}

func (r *ReserveRequest) Save(db *gorm.DB) {
	db.FirstOrCreate(&r)
}

func (r *ReserveRequest) Update(db *gorm.DB) {

}

func (r *ReserveRequest) Delete(db *gorm.DB) {
	db.Delete(r)
}
