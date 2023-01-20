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
		zap.L().Warn("User logging in failed")
	} else {
		zap.L().Info("User successfully logged in user")
	}
	defer res.Body.Close()
	return client
}

func CheckIsLoggedIn(client *http.Client) bool {
	// TODO: make checking equivalency with unlogged in text
	zap.L().Info("Started checking is user logged in")
	res, err := soup.GetWithClient(funcs.Linkefy("dateOfVisitDecision.do?siteLanguage="), client)
	if err != nil {
		zap.L().Warn("Can't get dateOfVisitDecision.do?siteLanguage=")
	}
	doc := soup.HTMLParse(res)
	td := doc.Find("td", "class", "infoTableInformationText")
	if td.Error != nil {
		zap.L().Warn("Element td doesnt exist on in the response")
	}
	zap.L().Info("Finished checking is user logged in")
	return funcs.StripString(td.Text()) == "There`s no reservation for this application."
}
