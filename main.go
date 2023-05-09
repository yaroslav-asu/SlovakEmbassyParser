package main

import (
	"main/internal/session"
	"main/internal/utils/db"
	"main/internal/utils/funcs"
	gorm_models "main/models/gorm"
)

func main() {
	funcs.Init()
	d := db.Connect()
	defer db.Close(d)
	var u gorm_models.User
	d.Model(&gorm_models.User{}).Where("id = 4").First(&u)
	s := session.NewLoggedInSession(u.UserName, u.Password)
	s.ParseCaptchas(100)
}
