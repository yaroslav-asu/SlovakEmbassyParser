package models

import (
	"gorm.io/gorm"
	"time"
)

type DayCell struct {
	gorm.Model
	AvailableReservations int
	CityId                string
	Date                  time.Time
}
