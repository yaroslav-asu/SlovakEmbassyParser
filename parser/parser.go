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
	client := session.LogIn()
	now := time.Now()
	return Parser{
		Session: client,
		Db:      db.Connect(),
		Date:    models.NewDateYM(now.Year(), int(now.Month())),
	}
}
func (p *Parser) Deconstruct() {
	zap.L().Info("Started parser deconstruction")
	session.LogOut(p.Session)
	zap.L().Info("Finished parser deconstruction")
}

func (p *Parser) getSoup(link string) (string, error) {
	return soup.GetWithClient(link, p.Session)
}

func (p *Parser) getParsedSoup(link string) soup.Root {
	doc, err := p.getSoup(link)
	if err != nil {
		zap.L().Error("Cant get: " + link)
		return soup.Root{}
	}
	return soup.HTMLParse(doc)
}

func (p *Parser) SaveToDB(model models.DbModel) {
	model.SaveToDB(p.Db)
}

func (p *Parser) DeleteFromDB(model models.DbModel) {
	model.DeleteFromDB(p.Db)
}
