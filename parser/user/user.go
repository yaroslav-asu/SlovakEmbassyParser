package user

import (
	"go.uber.org/zap"
	"main/internal/session"
	"main/internal/utils/db"
	"main/internal/utils/funcs"
	gorm_models "main/models/gorm"
	"main/models/gorm/datetime"
	"net/url"
	"strconv"
	"strings"
)

type User struct {
	Session session.Session
}

func (u *User) LogIn() {
	u.Session.LogInOnline()
}

func (u *User) LogOut() {
	zap.L().Info("Starting logout user: " + u.Session.User.UserName)
	u.Session.LogOut()
	zap.L().Info("Finished logout user: " + u.Session.User.UserName)
}

func NewUser(username, password string) User {
	zap.L().Info("Creating new user")
	newUser := User{
		Session: session.NewLoggedInSession(username, password),
	}
	return newUser
}

func NewUserFromModel(user gorm_models.User) User {
	return User{
		Session: session.NewLoggedInSession(user.UserName, user.Password),
	}
}

func (u *User) ReserveDatetime(city gorm_models.City, date datetime.Date) bool {
	zap.L().Info("Starting to reserve date in: " + city.Name + " at: " + date.Format(datetime.DateTime))
	res := u.Session.Get(funcs.Linkify("calendarDay.do?day=", date.Format(datetime.DateOnly), "&consularPostId=", city.Id))
	defer res.Body.Close()
	funcs.Sleep()
	captchaSolve := funcs.StripString(u.Session.SolveNewCaptcha())
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
	u.Session.User.IsReserved = u.IsReserved()
	if !u.Session.User.IsReserved {
		zap.L().Info("Failed to reserve user, trying again")
		for i := 0; i < 3; i++ {
			zap.L().Info(strconv.Itoa(i+1) + " try to reserve user")
			if u.ReserveDatetime(city, date) {
				u.deleteReserveRequests()
				return true
			}
		}
		return false
	} else {
		u.deleteReserveRequests()
	}
	return u.Session.User.IsReserved
}

func (u *User) IsReserved() bool {
	doc := u.Session.GetParsedSoup(funcs.Linkify("dateOfVisitDecision.do"))
	return strings.Contains(doc.Find("td", "class", "infoTableInformationText").Text(), "have reservation")
}

func (u *User) deleteReserveRequests() {
	dataBase := db.Connect()
	defer db.Close(dataBase)
	var reserveRequests []gorm_models.ReserveRequest
	err := dataBase.Model(gorm_models.ReserveRequest{}).Where("user_id = ?", u.Session.User.Id).Find(&reserveRequests).Error
	if err != nil {
		zap.L().Error("Failed to find user reserve requests: " + err.Error())
	}
	err = dataBase.Delete(reserveRequests).Error
	if err != nil {
		zap.L().Error("Failed to delete user requests: " + err.Error())
	}
}
