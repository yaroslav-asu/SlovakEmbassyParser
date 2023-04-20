package parser

import (
	"main/internal/session"
	"main/internal/utils/db"
	"main/internal/utils/vars"
	"main/models"
	gorm_models "main/models/gorm/datetime"
	"time"
)

type Parser struct {
	Session session.Session
	DB      *db.DB
	Date    gorm_models.Date
}

func NewParser() Parser {
	now := time.Now()
	return Parser{
		Session: session.NewLoggedInSession(vars.DefaultUserName, vars.DefaultUserPassword),
		DB:      db.Connect(),
		Date:    gorm_models.NewDateYM(now.Year(), int(now.Month())),
	}
}

func (p *Parser) SaveToDB(model models.DbModel) {
	model.SaveToDB(p.DB)
}

func (p *Parser) DeleteFromDB(model models.DbModel) {
	model.DeleteFromDB(p.DB)
}

func (p *Parser) Deconstruct() {
	p.DB.Close()
}
