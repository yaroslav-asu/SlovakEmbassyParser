package parser

import "main/internal/session"

func (p *Parser) Logout() {
	session.LogOut(p.Session)
}
