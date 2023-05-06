package models

import (
	"gorm.io/gorm"
)

type DBModel interface {
	Save(db *gorm.DB)
	Delete(db *gorm.DB)
}

type DBModelUpdater interface {
	Update(db *gorm.DB)
}
