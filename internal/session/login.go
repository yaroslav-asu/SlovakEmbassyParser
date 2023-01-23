package session

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"main/internal/utils/vars"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func Login() *http.Client {
	zap.L().Info("Started to log in user")
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		zap.L().Warn("Failed to create cookie jar")
	}
	client := &http.Client{Jar: cookieJar}
	_, err = client.Get(funcs.Linkefy("session.do"))
	if err != nil {
		zap.L().Warn("Can't get session.do for cookies")
	}
	res, err := client.PostForm(funcs.Linkefy("j_spring_security_check"), url.Values{"j_username": {vars.DefaultUserName}, "j_password": {vars.DefaultUserPassword}})
	if err != nil {
		zap.L().Warn("Can't post form to log in")
	}
	res, err = client.Get(funcs.Linkefy("dateOfVisitDecision.do?siteLanguage="))
	if err != nil {
		zap.L().Warn("Can't get dateOfVisitDecision.do?siteLanguage=")
	}
	if !CheckIsLoggedIn(client) {
		zap.L().Fatal("User logging in failed")
	} else {
		zap.L().Info("User successfully logged in")
	}
	defer res.Body.Close()
	return client
}

func CheckIsLoggedIn(client *http.Client) bool {
	zap.L().Info("Started checking is user logged in")
	loggedRes, err := soup.GetWithClient(funcs.Linkefy("dateOfVisitDecision.do?siteLanguage="), client)
	if err != nil {
		zap.L().Error("Got error while accessing to greeting page with session:\n" + err.Error())
	}
	loggedDoc := soup.HTMLParse(loggedRes)
	loggedText := loggedDoc.Find("table", "class", "infoTable").FullText()
	zap.L().Info("Got text with session")
	funcs.RandomSleep()
	unloggedRes, err := soup.Get(funcs.Linkefy("dateOfVisitDecision.do?siteLanguage="))
	if err != nil {
		zap.L().Error("Got error while accessing to greeting page without session:\n" + err.Error())
	}
	unloggedDoc := soup.HTMLParse(unloggedRes)
	unloggedText := unloggedDoc.Find("table", "class", "infoTable").FullText()
	zap.L().Info("Got text without session")
	zap.L().Info("Finished checking is user logged in")
	return loggedText != unloggedText
}
