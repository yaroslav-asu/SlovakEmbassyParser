package models

import (
	"gorm.io/gorm"
)

type DbModel interface {
	Save(db *gorm.DB)
	Update(db *gorm.DB)
	Delete(db *gorm.DB)
}
