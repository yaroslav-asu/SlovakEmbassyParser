package session

import (
	"net/http"
)

type Session struct {
	Client *http.Client
}

func NewSession() Session {
	return Session{
		Client: &http.Client{},
	}
}

func NewLoggedSession(username, password string) Session {
	session := Session{
		Client: &http.Client{},
	}
	session.LogIn(username, password)
	return session
}
