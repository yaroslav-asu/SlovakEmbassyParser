package gorm

import (
	"main/internal/utils/db"
)

type User struct {
	Id         uint `gorm:"primaryKey"`
	UserName   string
	Password   string
	IsReserved bool `gorm:"default:false"`
	TelegramId string
}

func (u *User) SaveToDB(db *db.DB) {
	db.DB.FirstOrCreate(u)
}
func (u *User) DeleteFromDB(db *db.DB) {
	db.DB.Delete(u)
}
