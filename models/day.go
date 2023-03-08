package models

import "gorm.io/gorm"

type DayCell struct {
	AvailableReservations int
	CityId                string
	Date                  Date
	gorm.Model
}
