package main

import (
	"main/internal/utils/funcs"
	"main/models"
	"main/parser"
)

func main() {
	funcs.Init()
	siteParser := parser.NewParser()
	defer siteParser.Deconstruct()
	date := models.NewDateNow()
	elems := siteParser.ParseAvailableReservations(models.City{Id: "542"}, date)
	siteParser.SaveToDB(elems)
}
