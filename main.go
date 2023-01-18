package main

import (
	"fmt"
	"main/internal/utils/funcs"
	"main/parser"
)

func main() {
	funcs.Init()
	siteParser := parser.NewParser()
	res, _ := siteParser.GetMonth(siteParser.WorkingCityByIndex(0), 1, 2023)
	fmt.Println(res)
}
