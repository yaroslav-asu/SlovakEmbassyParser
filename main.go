package main

import (
	"fmt"
	"main/internal/utils/funcs"
	"main/models"
	"main/parser"
)

func main() {
	funcs.Init()
	siteParser := parser.NewParser()
	parsedMonth := siteParser.ParseMonth(models.City{Id: "601"}, 3, 2023)
	fmt.Println(parsedMonth)
	parsedMonth = siteParser.ParseMonth(models.City{Id: "601"}, 4, 2023)
	fmt.Println(parsedMonth)
}
