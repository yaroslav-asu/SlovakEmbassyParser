package main

import (
	"main/internal/session"
	"main/internal/utils/funcs"
)

func main() {
	funcs.Init()
	s := session.NewLoggedInSession("herytlndten", "7753224")
	s.SolveNewCaptcha()
}
