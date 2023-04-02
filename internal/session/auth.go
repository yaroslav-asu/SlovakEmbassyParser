package session

import (
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"net/http/cookiejar"
	"net/url"
)

func (s *Session) LogIn(username, password string) {
	zap.L().Info("Started to log in user: " + username)
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		zap.L().Warn("Failed to create cookie jar")
	}
	s.Client.Jar = cookieJar
	s.Get(funcs.Linkify("session.do"))
	res := s.PostForm(
		funcs.Linkify("j_spring_security_check"),
		url.Values{
			"j_username": {username},
			"j_password": {password},
		},
	)
	res = s.Get(funcs.Linkify("dateOfVisitDecision.do?siteLanguage="))
	if !s.IsLoggedIn() {
		zap.L().Fatal("User login failed")
	} else {
		zap.L().Info("User successfully logged in")
	}
	defer res.Body.Close()
}

func (s *Session) IsLoggedIn() bool {
	zap.L().Info("Started checking is user logged in")
	loggedInSessionDoc := s.GetParsedSoup(funcs.Linkify("dateOfVisitDecision.do"))
	loggedInSessionText := loggedInSessionDoc.Find("table", "class", "infoTable").FullText()
	zap.L().Info("Got text with session")
	funcs.Sleep()
	blankSession := NewSession()
	blankSessionDoc := blankSession.GetParsedSoup(funcs.Linkify("dateOfVisitDecision.do"))
	blankSessionText := blankSessionDoc.Find("table", "class", "infoTable").FullText()
	zap.L().Info("Got text without session")
	defer zap.L().Info("Finished checking is user logged in")
	return loggedInSessionText != blankSessionText
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
