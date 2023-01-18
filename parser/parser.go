package parser

import (
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"time"
)

type Parser struct {
	client *http.Client
}

func NewParser(client *http.Client) Parser {
	return Parser{
		client: client,
	}
}

func (p Parser) RandomSleep() {
	sleepingTime := rand.Float64()*2 + 1
	zap.L().Info("Started random sleeping with time: " + fmt.Sprintf("%v", sleepingTime))
	time.Sleep(time.Duration(sleepingTime) * time.Second)
	zap.L().Info("Finished random sleeping")
}
