package main

import (
	"main/internal/logger"
	"main/internal/login"
	"main/internal/utils/db"
	"main/internal/utils/random"
	"main/internal/utils/variables"
	"main/parser"
)

func main() {
	variables.InitEnv()
	logger.InitLogger()
	random.InitRandom()
	client := login.Login()
	siteParser := parser.NewParser(client)
	dataBase := db.Connect()
	db.UpdateCities(dataBase, siteParser)
	//dataBase.Create(&models.City{Name: "asdf", CityId: "123"})

	//dataBase.First(&models.City{CityId: "123", Name: "asdf"}, "city_id = ?", "123")
	//siteParser.UpdateEmbassyCities(dataBase)
}
