package models

import "gorm.io/gorm"

type LastWorkingMonth struct {
	CityId string
	Date   Date
	gorm.Model
}
