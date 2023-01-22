package main

import (
	"main/internal/utils/funcs"
	"main/models"
	"main/parser"
	"time"
)

func main() {
	funcs.Init()
	siteParser := parser.NewParser()
	defer siteParser.Deconstruct()
	date, _ := time.Parse("02.01.2006", "02.02.2023")
	elems := siteParser.ParseAvailableReservations(models.City{Id: "542"}, date)
	siteParser.SaveToDB(elems)
}
