package parser

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/internal/session"
	"main/internal/utils/db"
	"main/internal/utils/funcs"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	session *http.Client
	Db      *gorm.DB
	Month   int
	Year    int
}

func NewParser() Parser {
	client := session.Login()
	now := time.Now()
	return Parser{
		session: client,
		Db:      db.Connect(),
		Month:   int(now.Month()),
		Year:    now.Year(),
	}
}

func (p *Parser) RandomSleep() {
	sleepingTime := rand.Float64()*2 + 1
	zap.L().Info("Started random sleeping with time: " + fmt.Sprintf("%v", sleepingTime))
	time.Sleep(time.Duration(sleepingTime) * time.Second)
	zap.L().Info("Finished random sleeping")
}

func (p *Parser) Get(link string) (string, error) {
	return soup.GetWithClient(link, p.session)
}

func ParseReservationData(data string) int {
	for _, s := range []string{"[", "]"} {
		data = strings.Replace(data, s, "", -1)
	}
	parsedNumbers := strings.Split(data, "/")
	reservedNum, err := strconv.Atoi(parsedNumbers[0])
	if err != nil {

	}
	totalNum, err := strconv.Atoi(parsedNumbers[1])
	if err != nil {

	}
	return totalNum - reservedNum
}

func ParseMonthCellDate(date string, year int) time.Time {
	parsedDate := strings.Split(date, ".")
	intDate := funcs.StringsToIntArray(parsedDate)
	day, month := intDate[0], intDate[1]
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
