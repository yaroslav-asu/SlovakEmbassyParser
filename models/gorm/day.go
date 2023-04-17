package gorm

import (
	"main/models/gorm/datetime"
)

type DayCell struct {
	AvailableReservations int
	CityId                string
	City                  City
	Date                  datetime.Date
}
