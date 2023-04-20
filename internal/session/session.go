package session

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/session/captcha"
	gorm_models "main/models/gorm"
	"net/http"
	"time"
)

const requestTimeout = 30 * time.Second

type Session struct {
	Client   *http.Client
	Proxy    gorm_models.Proxy
	captcha  captcha.Captcha
	username string
	password string
}

func NewBlankSession() Session {
	return Session{
		Client: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func NewBlankProxiedSession() Session {
	session := NewBlankSession()
	session.ChangeProxy()
	return session
}

func NewSession(username, password string) Session {
	session := NewBlankProxiedSession()
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
