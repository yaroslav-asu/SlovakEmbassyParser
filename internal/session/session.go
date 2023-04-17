package session

import (
	gorm_models "main/models/gorm"
	"net/http"
	"time"
)

const requestTimeout = 30 * time.Second

type Session struct {
	Client *http.Client
	Proxy  gorm_models.Proxy
}

func NewSession() Session {
	session := Session{
		Client: &http.Client{
			Timeout: requestTimeout,
		},
	}
	session.ChangeProxy()
	return session
}

func NewLoggedInSession(username, password string) Session {
	session := NewSession()
	session.LogIn(username, password)
	return session
}
