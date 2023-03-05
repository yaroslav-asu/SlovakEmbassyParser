package models

import "gorm.io/gorm"

type DbModel interface {
	SaveToDb(db *gorm.DB)
}
