package session

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	gorm_models "main/models/gorm"
	"net/http"
	"time"
)

const requestTimeout = 30 * time.Second

type Session struct {
	Client   *http.Client
	Proxy    gorm_models.Proxy
	username string
	password string
}

func NewBlankSession() Session {
	session := Session{
		Client: &http.Client{
			Timeout: requestTimeout,
		},
	}
	session.ChangeProxy()
	return session
}

func NewSession(username, password string) Session {
	session := NewBlankSession()
	session.username = username
	session.password = password
	return session
}

func NewLoggedInSession(username, password string) Session {
	session := NewSession(username, password)
	session.LogIn()
	return session
}

func sessionWorking(root soup.Root) bool {
	if len(root.FindAll("input", "id", "j_username")) > 0 || len(root.FindAll("input", "id", "j_password")) > 0 {
		zap.L().Warn("Session expired")
		return false
	}
	zap.L().Info("Session still work")
	return true
}
