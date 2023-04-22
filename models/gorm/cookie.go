package gorm

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type Cookie struct {
	Id     uint `gorm:"primaryKey"`
	UserId uint
	User   User

	Name  string
	Value string
}

func NewCookie(owner User, cookie http.Cookie) Cookie {
	return Cookie{
		UserId: owner.Id,
		User:   owner,
		Name:   cookie.Name,
		Value:  cookie.Value,
	}
}

func (c Cookie) Cookie() *http.Cookie {
	return &http.Cookie{
		Name:  c.Name,
		Value: c.Value,
	}
}

func (c Cookie) SaveToDB(db *gorm.DB) {
	zap.L().Info("Saved to DB")
	db.FirstOrCreate(&c, c)
}

func (c Cookie) DeleteFromDB(db *gorm.DB) {
	db.Where("id = ?", c.Id).Delete(&c)
}

func (c Cookie) Format() string {
	return fmt.Sprintf("Cookie{Id: %d, Name: %s, Value: %s}", c.Id, c.Name, c.Value)
}
