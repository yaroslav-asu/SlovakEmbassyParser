package main

import (
	"fmt"
	"main/internal/utils/funcs"
	"main/models/gorm/datetime"
	"main/parser"
)

func main() {
	funcs.Init()
	p := parser.NewLoggedInParser()
	date, err := p.CurrentMonth()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(date.Format(datetime.DateTime))
}
