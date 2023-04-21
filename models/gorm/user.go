package gorm

import "gorm.io/gorm"

type User struct {
	Id         uint `gorm:"primaryKey"`
	UserName   string
	Password   string
	IsReserved bool `gorm:"default:false"`
	TelegramId string
}

func (u *User) SaveToDB(db *gorm.DB) {
	db.FirstOrCreate(u)
}
func (u *User) DeleteFromDB(db *gorm.DB) {
	db.Delete(u)
}
