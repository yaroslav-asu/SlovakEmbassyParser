package main

import (
	"main/internal/logger"
	"main/internal/login"
	"main/internal/variables"
)

func main() {
	variables.InitEnv()
	logger.InitLogger()

	_ = login.Login()
	//fmt.Println(login.CheckIsLoggedIn(client))
}
