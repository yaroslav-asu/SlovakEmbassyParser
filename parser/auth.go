package parser

func (p *Parser) LogIn() {
	p.Session.LogInOnline()
}

func (p *Parser) LogOut() {
	p.Session.LogOut()
}
