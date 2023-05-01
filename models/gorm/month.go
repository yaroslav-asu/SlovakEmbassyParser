package gorm

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/models/gorm/datetime"
)

type Month struct {
	Id   uint `gorm:"primaryKey"`
	Date datetime.Date
}

func NewMonth(date datetime.Date) Month {
	return Month{
		Date: date,
	}
}

func (m Month) Save(db *gorm.DB) {
	zap.L().Info("Month " + m.Date.Format(datetime.MonthAndYear) + " saved to DB")
	db.FirstOrCreate(&m, m)
}

func (m Month) Delete(db *gorm.DB) {
	db.Where("id = ?", m.Id).Delete(&m)
}

func (m Month) Format() string {
	return fmt.Sprintf("Month{DateTime: %s}", m.Date.Format(datetime.MonthAndYear))
}
