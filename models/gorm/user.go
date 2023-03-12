package gorm

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"main/internal/datetime"
	"main/internal/session"
	"main/internal/utils/funcs"
	"net/http"
	"net/url"
	"os"
)

type User struct {
	Id         uint `gorm:"primaryKey"`
	UserName   string
	Password   string
	TelegramId string
	Session    *http.Client `gorm:"-:migration"`
}

func (u *User) Login() {
	u.Session = session.LogIn(u.UserName, u.Password)
}

func NewUser(username, password string) User {
	newUser := User{
		UserName: username,
		Password: password,
	}
	newUser.Login()
	return newUser
}

func (u *User) DownloadCaptcha() {
	res, err := u.Session.Get(funcs.Linkify("simpleCaptcha.png"))
	if err != nil {
		zap.L().Error(err.Error())
	}
	defer res.Body.Close()
	file, err := os.Create("captcha.png")
	if err != nil {
		zap.L().Error("Cant create captcha.png with error: " + err.Error())
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		zap.L().Error("Cant write captcha bytes to file with error: " + err.Error())
	}
}

func (u *User) SolveCaptcha() string {
	var textCaptcha string
	fmt.Print("Type captcha solve: ")
	_, err := fmt.Scan(&textCaptcha)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return textCaptcha
}

func (u *User) ReserveDatetime(user User, city City, date datetime.Date) {
	zap.L().Info("Starting to reserve date in: " + city.Name + " at: " + date.Format(datetime.DateTime))
	res, err := user.Session.Get(funcs.Linkify("calendarDay.do?day=", date.Format(datetime.DateOnly), "&timeSlotId=&calendarId=&consularPostId=", city.Id))
	defer res.Body.Close()
	if err != nil {
		zap.L().Error("Cant get page of reserving date with error: " + err.Error())
	}
	funcs.Sleep()
	user.DownloadCaptcha()
	captchaSolve := user.SolveCaptcha()
	res, err = user.Session.PostForm(
		funcs.Linkify("calendarDay.do?day=", date.Format(datetime.DateOnly), "&timeSlotId=&calendarId=&consularPostId=", city.Id),
		url.Values{
			"calendar.timeOfVisit":               {date.Format(datetime.FormDateTime)},
			"calendar.sequenceNo":                {"1"},
			"calendar.consularPost.consularPost": {city.Id},
			"captcha":                            {captchaSolve},
			"calendar.timeSlot.timeSlotId":       {""},
			"calendar.calendarId":                {""},
		},
	)
	defer res.Body.Close()
	if err != nil {
		zap.L().Error("Cant post reserve form with error: " + err.Error())
	}
	res, err = user.Session.Get(funcs.Linkify("logout.do"))
	defer res.Body.Close()
	if err != nil {
		zap.L().Error("Cant get logout.do with error: " + err.Error())
	}
	funcs.Sleep()
	session.LogOut(user.Session)
}
