package parser

import (
	"main/internal/utils/vars"
)

func (p *Parser) LogIn() {
	p.Session.LogIn(vars.DefaultUserName, vars.DefaultUserPassword)
}

func (p *Parser) LogOut() {
	p.Session.LogOut()
}
