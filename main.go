package main

import (
	"fmt"
	"main/internal/session"
	"main/internal/utils/funcs"
	"main/models"
	"main/parser"
	"net/http"
)

func main() {
	funcs.Init()
	siteParser := parser.NewParser()
	session.Logout(siteParser.Session)
	fmt.Println(siteParser.ParseMonth(models.City{Id: "601"}, 1, 2023))
	res, _ := http.Get(funcs.Linkefy("calendar.do?month=0&consularPost=577"))
	fmt.Println(res.StatusCode)
	//parsedMonth := siteParser.ParseMonth(models.City{Id: "601"}, 3, 2023)
	//fmt.Println(parsedMonth)
	//parsedMonth = siteParser.ParseMonth(models.City{Id: "601"}, 4, 2023)
	//fmt.Println(parsedMonth)
}
