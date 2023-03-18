package gorm

import (
	"gorm.io/gorm"
	"main/models/gorm/datetime"
)

type DayCell struct {
	AvailableReservations int
	CityId                string
	Date                  datetime.Date
	gorm.Model
}
