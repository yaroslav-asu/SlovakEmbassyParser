package main

import (
	"fmt"
	"main/internal/session"
	"main/internal/utils/funcs"
)

func main() {
	funcs.Init()
	s := session.NewBlankProxiedSession()
	solve := s.SolveNewCaptcha()
	fmt.Println(solve)
}
