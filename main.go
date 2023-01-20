package main

import (
	"fmt"
	"main/internal/session"
	"main/internal/utils/funcs"
	"main/parser"
)

func main() {
	funcs.Init()
	siteParser := parser.NewParser()
	defer siteParser.Deconstruct()
	fmt.Println(session.CheckIsLoggedIn(siteParser.Session))
}
