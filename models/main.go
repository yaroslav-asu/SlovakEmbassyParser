package models

import "gorm.io/gorm"

type DbModel interface {
	FirstOrCreate(db *gorm.DB)
}

type DbModelArray interface {
	SaveToDb(db *gorm.DB)
}
