package main

import (
	"fmt"
	"main/internal/logger"
	"main/internal/login"
	"main/internal/utils/random"
	"main/internal/utils/variables"
	"main/parser"
)

func main() {
	variables.InitEnv()
	logger.InitLogger()
	random.InitRandom()

	client := login.Login()
	siteParser := parser.NewParser(client)
	fmt.Println(siteParser.GetEmbassyCities())
}
