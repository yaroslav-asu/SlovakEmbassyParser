package parser

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/internal/session"
	"main/internal/utils/db"
	"main/models"
	"net/http"
	"time"
)

type Parser struct {
	Session *http.Client
	Db      *gorm.DB
	Date    models.Date
}

func NewParser() Parser {
	client := session.Login()
	now := time.Now()
	return Parser{
		Session: client,
		Db:      db.Connect(),
		Date:    models.NewDateYM(now.Year(), int(now.Month())),
	}
}
func (p *Parser) Deconstruct() {
	zap.L().Info("Started parser deconstruction")
	session.Logout(p.Session)
	zap.L().Info("Finished parser deconstruction")
}

func (p *Parser) GetSoup(link string) (string, error) {
	return soup.GetWithClient(link, p.Session)
}

func (p *Parser) SaveToDB(DatabaseModel models.DbModel) {
	DatabaseModel.SaveToDb(p.Db)
}
