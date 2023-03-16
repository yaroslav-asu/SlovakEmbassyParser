package gorm

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/internal/datetime"
	"main/internal/session"
	"main/internal/utils/funcs"
	"net/url"
	"strings"
)

type User struct {
	Id         uint `gorm:"primaryKey"`
	UserName   string
	Password   string
	TelegramId string
	Session    session.Session `gorm:"-:all"`
}

func (u *User) SaveToDB(db *gorm.DB) {
	db.FirstOrCreate(u)
}
func (u *User) DeleteFromDB(db *gorm.DB) {

}
func (u *User) LogIn() {
	u.Session.LogIn(u.UserName, u.Password)
}

func (u *User) LogOut() {
	zap.L().Info("Starting logout user: " + u.UserName)
	u.Session.LogOut()
	zap.L().Info("Finished logout user: " + u.UserName)
}

func NewUser(username, password string) User {
	newUser := User{
		UserName: username,
		Password: password,
		Session:  session.NewLoggedSession(username, password),
	}
	return newUser
}

func (u *User) ReserveDatetime(city City, date datetime.Date) {
	zap.L().Info("Starting to reserve date in: " + city.Name + " at: " + date.Format(datetime.DateTime))
	res := u.Session.Get(funcs.Linkify("calendarDay.do?day=", date.Format(datetime.DateOnly), "&timeSlotId=&calendarId=&consularPostId=", city.Id))
	defer res.Body.Close()
	funcs.Sleep()
	u.Session.DownloadCaptcha()
	captchaSolve := session.SolveCaptcha()
	res = u.Session.PostForm(
		funcs.Linkify("calendarDay.do?day=", date.Format(datetime.DateOnly), "&consularPostId=", city.Id),
		url.Values{
			"calendar.timeOfVisit":               {date.Format(datetime.FormDateTime)},
			"calendar.sequenceNo":                {"1"},
			"calendar.consularPost.consularPost": {city.Id},
			"captcha":                            {captchaSolve},
		},
	)
	defer res.Body.Close()
	res = u.Session.Get(funcs.Linkify("logout.do"))
	defer res.Body.Close()
}

func (u *User) IsReserved() bool {
	doc := u.Session.GetParsedSoup(funcs.Linkify("dateOfVisitDecision.do"))
	return strings.Contains(doc.Find("td", "class", "infoTableInformationText").Text(), "have reservation")
}
