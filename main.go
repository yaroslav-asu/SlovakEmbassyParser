package main

import (
	"main/internal/login"
	"main/internal/utils/db"
	"main/internal/utils/init"
	"main/parser"
)

func main() {
	init.Init()
	client := login.Login()
	siteParser := parser.NewParser(client)
	dataBase := db.Connect()
	db.UpdateCities(dataBase, siteParser)
}
