package parser

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/internal/login"
	"main/internal/utils/db"
	"math/rand"
	"net/http"
	"time"
)

type Parser struct {
	session *http.Client
	db      *gorm.DB
	month   int
	year    int
}

func NewParser() Parser {
	client := login.Login()
	now := time.Now()
	return Parser{
		session: client,
		db:      db.Connect(),
		month:   int(now.Month()),
		year:    now.Year(),
	}
}

func (p Parser) RandomSleep() {
	sleepingTime := rand.Float64()*2 + 1
	zap.L().Info("Started random sleeping with time: " + fmt.Sprintf("%v", sleepingTime))
	time.Sleep(time.Duration(sleepingTime) * time.Second)
	zap.L().Info("Finished random sleeping")
}

func (p Parser) Get(link string) (string, error) {
	return soup.GetWithClient(link, p.session)
}
