package session

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/session/captcha"
	"main/internal/utils/db"
	gorm_models "main/models/gorm"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const requestTimeout = 15 * time.Second

type Session struct {
	Client  *http.Client
	Proxy   gorm_models.Proxy
	captcha captcha.Captcha
	User    gorm_models.User
}

func NewBlankSession() Session {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		zap.L().Error("Failed to create cookie jar")
	}
	return Session{
		Client: &http.Client{
			Timeout: requestTimeout,
			Jar:     cookieJar,
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
	dataBase := db.Connect()
	defer db.Close(dataBase)
	user := gorm_models.User{
		UserName: username,
		Password: password,
	}
	dataBase.Find(&user)
	session.User = user
	return session
}

func NewLoggedInSession(username, password string) Session {
	session := NewSession(username, password)
	session.LogIn()
	return session
}

func isLoggedIn(root soup.Root) bool {
	if len(root.FindAll("input", "id", "j_username")) > 0 || len(root.FindAll("input", "id", "j_password")) > 0 {
		zap.L().Warn("Session expired")
		return false
	}
	zap.L().Info("Session still work")
	return true
}
