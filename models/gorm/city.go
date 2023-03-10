package gorm

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/internal/datetime"
)

type City struct {
	Id           string
	Name         string
	StartWorking datetime.Date
	EndWorking   datetime.Date
}

func (c City) SaveToDB(db *gorm.DB) {
	zap.L().Info("Saved to DB")
	db.FirstOrCreate(&c, c)
}

func (c City) DeleteFromDB(db *gorm.DB) {
	db.Where("id = ? and name = ?", c.Id, c.Name).Delete(&c)
}
