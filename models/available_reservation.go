package models

import (
	"gorm.io/gorm"
	"time"
)

type AvailableReservation struct {
	gorm.Model
	Date   time.Time
	CityId string
}
