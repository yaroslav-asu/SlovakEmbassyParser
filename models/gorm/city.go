package gorm

import (
	"go.uber.org/zap"
	"main/internal/utils/db"
	"main/models/gorm/datetime"
)

type City struct {
	Id           string
	Name         string
	StartWorking datetime.Date
	EndWorking   datetime.Date
}

func (c City) SaveToDB(db *db.DB) {
	zap.L().Info("Saved to DB")
	db.DB.FirstOrCreate(&c, c)
}

func (c City) DeleteFromDB(db *db.DB) {
	db.DB.Where("id = ? and name = ?", c.Id, c.Name).Delete(&c)
}
