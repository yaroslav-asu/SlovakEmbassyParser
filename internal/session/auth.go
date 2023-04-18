package session

import (
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"net/http/cookiejar"
	"net/url"
)

func (s *Session) LogIn() {
	zap.L().Info("Started to log in user: " + s.username)
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		zap.L().Error("Failed to create cookie jar")
	}
	s.Client.Jar = cookieJar
	s.Get(funcs.Linkify("session.do"))
	res := s.PostForm(
		funcs.Linkify("j_spring_security_check"),
		url.Values{
			"j_username": {s.username},
			"j_password": {s.password},
		},
	)
	defer res.Body.Close()
	res = s.Get(funcs.Linkify("dateOfVisitDecision.do?siteLanguage="))
	root := funcs.ResponseToSoup(res)
	if !sessionWorking(root) {
		zap.L().Fatal("User login failed")
	} else {
		zap.L().Info("User successfully logged in")
	}
	defer res.Body.Close()
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
