package session

import (
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"main/internal/utils/db"
	"main/internal/utils/vars"
	gorm_models "main/models/gorm"
	"net/url"
)

func (s *Session) SaveCookiesToDb() {
	zap.L().Info("Saving cookies to db")
	siteUrl, err := url.Parse(vars.SiteUrl)
	if err != nil {
		zap.L().Error("Failed to parse site url: " + vars.SiteUrl)
		return
	}
	dataBase := db.Connect()
	defer db.Close(dataBase)
	for _, cookie := range s.Client.Jar.Cookies(siteUrl) {
		gorm_models.NewCookie(s.User, *cookie).SaveOrCreate(dataBase)
	}
}

func (s *Session) savedSessionCookies() []gorm_models.Cookie {
	zap.L().Info("Getting saved cookies to user: " + s.User.UserName + " from db")
	dataBase := db.Connect()
	defer db.Close(dataBase)
	var cookies []gorm_models.Cookie
	dataBase.Preload(clause.Associations).Find(&cookies)
	dataBase.Model(&gorm_models.Cookie{}).Joins(" join users on cookies.user_id = users.id").Where("user_name = ?", s.User.UserName).Find(&cookies)
	zap.L().Info("Got session cookies")
	return cookies
}
