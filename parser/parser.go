package parser

import (
	"gorm.io/gorm"
	"main/internal/session"
	"main/internal/utils/db"
	"main/internal/utils/vars"
	"main/models"
	gorm_models "main/models/gorm/datetime"
	"time"
)

type Parser struct {
	Session session.Session
	DB      *gorm.DB
	Date    gorm_models.Date
}

func NewLoggedInParser() Parser {
	parser := NewParser()
	parser.Session.LogIn()
	return parser
}
func NewParser() Parser {
	now := time.Now()
	return Parser{
		Session: session.NewSession(vars.DefaultUserName, vars.DefaultUserPassword),
		DB:      db.Connect(),
		Date:    gorm_models.NewDateYM(now.Year(), int(now.Month())),
	}
}

func (p *Parser) SaveToDB(model models.DBModel) {
	model.Save(p.DB)
}

func (p *Parser) DeleteFromDB(model models.DBModel) {
	model.Delete(p.DB)
}

func (p *Parser) Deconstruct() {
	db.Close(p.DB)
}
