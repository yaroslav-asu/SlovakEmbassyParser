package gorm

import (
	"fmt"
	"gorm.io/gorm"
	"main/models/gorm/datetime"
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

func (r *ReserveRequest) Format() string {
	return fmt.Sprintf("ReserveRequest{UserId: %d, User: %v, CityId: %s, City: %v, MonthId: %d, Month: %v}",
		r.UserId,
		r.User.Format(),
		r.CityId,
		r.City.Format(),
		r.MonthId,
		r.Month.Date.Format(datetime.MonthAndYear),
	)
}
