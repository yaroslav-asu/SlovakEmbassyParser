package parser

import (
	"gorm.io/gorm"
	"main/internal/session"
	"main/internal/utils/db"
	"main/internal/utils/vars"
	"main/models"
	"main/models/gorm/datetime"
)

type Parser struct {
	Session session.Session
	DB      *gorm.DB
	Date    datetime.Date
}

func NewLoggedInParser() Parser {
	parser := NewParser()
	parser.Session.LogIn()
	return parser
}
func NewParser() Parser {
	return Parser{
		Session: session.NewSession(vars.DefaultUserName, vars.DefaultUserPassword),
		DB:      db.Connect(),
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
