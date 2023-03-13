package main

import (
	"main/internal/datetime"
	"main/internal/utils/funcs"
	"main/models/gorm"
	"main/parser"
)

func main() {
	funcs.Init()
	siteParser := parser.NewParser()
	defer siteParser.LogOut()
	user := gorm.NewUser("herytlndten", "7753224")
	defer user.LogOut()
	user.ReserveDatetime(gorm.City{Id: "542", Name: "Abuja"}, datetime.NewDate(2023, 3, 30, 14, 30))
}
