package models

import "gorm.io/gorm"

type DbModel interface {
	SaveToDB(db *gorm.DB)
	DeleteFromDB(db *gorm.DB)
}
