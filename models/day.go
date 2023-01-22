package models

import (
	"gorm.io/gorm"
)

type DayCell struct {
	gorm.Model
	AvailableReservations int
	CityId                string
	Date                  Date
}
