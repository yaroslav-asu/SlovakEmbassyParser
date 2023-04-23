package gorm

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/models/gorm/datetime"
)

type City struct {
	Id           string
	Name         string
	StartWorking datetime.Date
	EndWorking   datetime.Date
}

func (c City) Update(db *gorm.DB) {

}
func (c City) Save(db *gorm.DB) {
	zap.L().Info("Saved to DB")
	db.FirstOrCreate(&c, c)
}

func (c City) Delete(db *gorm.DB) {
	db.Where("id = ? and name = ?", c.Id, c.Name).Delete(&c)
}

func (c City) Format() string {
	return fmt.Sprintf("City{Id: %s, Name: %s}", c.Id, c.Name)
}
