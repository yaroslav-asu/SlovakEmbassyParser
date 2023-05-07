package session

import (
	"errors"
	"go.uber.org/zap"
	"main/internal/utils/db"
	"main/internal/utils/funcs"
	"main/internal/utils/vars"
	"main/models/gorm/datetime"
	"net/http"
	"net/url"
)

var cookiesLogInError = errors.New("failed to log in by cookies")
var currentSessionDateParseError = errors.New("failed to parse current session date")

func (s *Session) LogIn() {
	var err error
	if len(s.savedSessionCookies()) > 0 {
		err = s.LogInWithCookies()
	} else {
		err = s.LogInOnline()
	}
	for err != nil {
		zap.L().Info("Failed to log in with error: " + err.Error() + " , trying to log in online")
		err = s.LogInOnline()
	}
	zap.L().Info("Successfully logged in")

}

func (s *Session) LogInOnline() error {
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
	root, err := funcs.ResponseToSoup(res)
	if err != nil {
		return err
	}
	if !isLoggedIn(root) {
		zap.L().Fatal("User login failed")
	} else {
		zap.L().Info("User successfully logged in")
	}
	s.SaveCookiesToDb()
	return nil
}

func (s *Session) LogInWithCookies() error {
	zap.L().Info("Started to log in user: " + s.User.UserName + " with saved cookies")
	dataBase := db.Connect()
	defer db.Close(dataBase)
	siteUrl, err := url.Parse(vars.SiteUrl)
	if err != nil {
		zap.L().Error("Failed to parse site url")
		return err
	}
	cookies := s.savedSessionCookies()
	for _, cookie := range cookies {
		s.Client.Jar.SetCookies(siteUrl, []*http.Cookie{cookie.Cookie()})
	}
	res := s.Get(funcs.Linkify("dateOfVisitDecision.do?siteLanguage="))
	defer res.Body.Close()
	root, err := funcs.ResponseToSoup(res)
	if err != nil {
		return err
	}
	if isLoggedIn(root) {
		zap.L().Info("Successfully logged in by cookies")
		var currentDate datetime.Date
		i := 0
		for err := errors.New(""); err != nil && i < 3; i++ {
			currentDate, err = s.GetDate()
			if err == nil {
				s.Date = currentDate
				zap.L().Info("Successfully got current session date: " + currentDate.Format(datetime.MonthAndYear))
				break
			}
		}
		if err != nil {
			zap.L().Info(currentSessionDateParseError.Error())
			return currentSessionDateParseError
		}
	} else {
		zap.L().Error(cookiesLogInError.Error())
		return cookiesLogInError
	}
	return nil
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
