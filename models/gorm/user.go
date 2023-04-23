package gorm

import "gorm.io/gorm"

type User struct {
	Id         uint `gorm:"primaryKey"`
	UserName   string
	Password   string
	IsReserved bool `gorm:"default:false"`
	TelegramId string
}

func (u *User) Save(db *gorm.DB) {
	db.FirstOrCreate(u)
}

func (u *User) Update(db *gorm.DB) {

}

func (u *User) Delete(db *gorm.DB) {
	db.Delete(u)
}
