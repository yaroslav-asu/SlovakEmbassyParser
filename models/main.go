package models

import "gorm.io/gorm"

type DbModel interface {
	FirstOrCreate(db *gorm.DB)
}

type SaveAble interface {
	SaveToDb(db *gorm.DB)
}
