package session

import (
	"go.uber.org/zap"
	"main/internal/utils/db"
	"main/internal/utils/funcs"
	"main/internal/utils/vars"
	"net/http"
	"net/url"
)

func (s *Session) LogIn() {
	if len(s.savedSessionCookies()) > 0 {
		s.LogInWithCookies()
	} else {
		s.LogInOnline()
	}
}

func (s *Session) LogInOnline() {
	zap.L().Info("Started to online log in user: " + s.User.UserName)
	s.Get(funcs.Linkify("session.do"))
	res := s.PostForm(
		funcs.Linkify("j_spring_security_check"),
		url.Values{
			"j_username": {s.User.UserName},
			"j_password": {s.User.Password},
		},
	)
	defer res.Body.Close()
	res = s.Get(funcs.Linkify("dateOfVisitDecision.do?siteLanguage="))
	defer res.Body.Close()
	root := funcs.ResponseToSoup(res)
	if !isLoggedIn(root) {
		zap.L().Fatal("User login failed")
	} else {
		zap.L().Info("User successfully logged in")
	}
	s.SaveCookiesToDb()
}

func (s *Session) LogInWithCookies() {
	zap.L().Info("Started to log in user: " + s.User.UserName + " with saved cookies")
	dataBase := db.Connect()
	defer db.Close(dataBase)
	siteUrl, err := url.Parse(vars.SiteUrl)
	if err != nil {
		zap.L().Error("Failed to parse site url")
		return
	}
	cookies := s.savedSessionCookies()
	for _, cookie := range cookies {
		s.Client.Jar.SetCookies(siteUrl, []*http.Cookie{cookie.Cookie()})
	}
	res := s.Get(funcs.Linkify("dateOfVisitDecision.do?siteLanguage="))
	defer res.Body.Close()
	root := funcs.ResponseToSoup(res)
	if isLoggedIn(root) {
		zap.L().Info("Successfully logged in by cookies")
	} else {
		zap.L().Error("Failed to login by cookies")
	}
}

func (s *Session) LogOut() {
	zap.L().Info("Starting to logout")
	res := s.Get(funcs.Linkify("j_spring_security_logout"))
	switch res.StatusCode {
	case 200:
		zap.L().Info("Successfully logged out")
	default:
		zap.L().Warn("On logout got error with code: " + res.Status)
	}
}
