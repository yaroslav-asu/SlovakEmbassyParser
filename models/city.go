package models

import (
	"gorm.io/gorm"
)

type City struct {
	Id      string
	Name    string
	Working bool
	gorm.Model
}
