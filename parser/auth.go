package parser

func (p *Parser) LogIn() {
	p.Session.LogIn()
}

func (p *Parser) LogOut() {
	p.Session.LogOut()
}
