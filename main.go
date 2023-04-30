package main

import (
	"main/internal/utils/funcs"
	"main/parser"
)

func main() {
	funcs.Init()
	p := parser.NewParser()
	p.RunCheckingReserveRequests()
}
