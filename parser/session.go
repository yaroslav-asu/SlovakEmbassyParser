package parser

import "main/internal/session"

func (p *Parser) Logout() {
	session.Logout(p.Session)
}
