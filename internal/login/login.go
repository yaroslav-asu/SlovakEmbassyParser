package login

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/functions"
	"main/internal/utils/variables"
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
	_, err = client.Get("https://ezov.mzv.sk/e-zov/login.do")
	if err != nil {
		zap.L().Warn("Can't get login.do for cookies")
	}
	res, err := client.PostForm("https://ezov.mzv.sk/e-zov/j_spring_security_check", url.Values{"j_username": {variables.DefaultUserName}, "j_password": {variables.DefaultUserPassword}})
	if err != nil {
		zap.L().Warn("Can't post form to log in")
	}
	res, err = client.Get("https://ezov.mzv.sk/e-zov/dateOfVisitDecision.do?siteLanguage=")
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
	zap.L().Info("Started checking is user logged in")
	res, err := soup.GetWithClient(variables.SiteUrl+"dateOfVisitDecision.do?siteLanguage=", client)
	if err != nil {
		zap.L().Warn("Can't get dateOfVisitDecision.do?siteLanguage=")
	}
	doc := soup.HTMLParse(res)
	td := doc.Find("td", "class", "infoTableInformationText")
	if td.Error != nil {
		zap.L().Warn("Element td doesnt exist on in the response")
	}
	zap.L().Info("Finished checking is user logged in")
	return functions.StripString(td.Text()) == "There`s no reservation for this application."
}
