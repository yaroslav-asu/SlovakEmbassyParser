package main

import (
	"main/internal/utils/funcs"
	"main/parser"
)

func main() {
	funcs.Init()
	siteParser := parser.NewParser()
	siteParser.RunCheckingReserveRequests()
	defer siteParser.Deconstruct()
}
