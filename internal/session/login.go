package session

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func LogIn(username, password string) *http.Client {
	zap.L().Info("Started to log in user: " + username)
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		zap.L().Warn("Failed to create cookie jar")
	}
	client := &http.Client{Jar: cookieJar}
	_, err = client.Get(funcs.Linkify("session.do"))
	if err != nil {
		zap.L().Warn("Can't get session.do cookies page")
	}
	res, err := client.PostForm(funcs.Linkify("j_spring_security_check"), url.Values{"j_username": {username}, "j_password": {password}})
	if err != nil {
		zap.L().Warn("Can't post form to log in")
	}
	res, err = client.Get(funcs.Linkify("dateOfVisitDecision.do?siteLanguage="))
	if err != nil {
		zap.L().Warn("Can't get dateOfVisitDecision.do?siteLanguage=")
	}
	if !IsLoggedIn(client) {
		zap.L().Fatal("User login failed")
	} else {
		zap.L().Info("User successfully logged in")
	}
	defer res.Body.Close()
	return client
}

func IsLoggedIn(client *http.Client) bool {
	zap.L().Info("Started checking is user logged in")
	loggedInRes, err := soup.GetWithClient(funcs.Linkify("dateOfVisitDecision.do?siteLanguage="), client)
	if err != nil {
		zap.L().Error("Got error while accessing to greeting page from session:\n" + err.Error())
	}
	loggedInDoc := soup.HTMLParse(loggedInRes)
	loggedText := loggedInDoc.Find("table", "class", "infoTable").FullText()
	zap.L().Info("Got text with session")
	funcs.Sleep()
	loggedOutRes, err := soup.Get(funcs.Linkify("dateOfVisitDecision.do?siteLanguage="))
	if err != nil {
		zap.L().Error("Got error while accessing to greeting page without session:\n" + err.Error())
	}
	loggedOutDoc := soup.HTMLParse(loggedOutRes)
	loggedOutText := loggedOutDoc.Find("table", "class", "infoTable").FullText()
	zap.L().Info("Got text without session")
	zap.L().Info("Finished checking is user logged in")
	return loggedText != loggedOutText
}
