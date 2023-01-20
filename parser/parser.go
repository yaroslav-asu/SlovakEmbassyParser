package parser

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/internal/session"
	"main/internal/utils/db"
	"main/internal/utils/funcs"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	Session *http.Client
	Db      *gorm.DB
	Month   int
	Year    int
}

func NewParser() Parser {
	client := session.Login()
	now := time.Now()
	return Parser{
		Session: client,
		Db:      db.Connect(),
		Month:   int(now.Month()),
		Year:    now.Year(),
	}
}
func (p *Parser) Deconstruct() {
	zap.L().Info("Started parser deconstruction")
	session.Logout(p.Session)
	zap.L().Info("Finished parser deconstruction")
}

func (p *Parser) RandomSleep() {
	funcs.RandomSleep()
}

func (p *Parser) Get(link string) (string, error) {
	return soup.GetWithClient(link, p.Session)
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
