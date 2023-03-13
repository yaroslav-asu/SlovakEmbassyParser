package session

import "net/http"

type Session struct {
	client http.Client
}

func NewSession(username, password string) Session {
	session := Session{
		client: http.Client{},
	}
	session.LogIn(username, password)
	return session
}
