package session

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	gorm_models "main/models/gorm"
	"main/models/gorm/datetime"
	"math"
	"strconv"
	"strings"
)

func monthDate(root soup.Root) (datetime.Date, error) {
	monthDate := root.Find("td", "class", "calendarMonthLabel")
	if monthDate.Error != nil {
		return datetime.Date{}, monthDate.Error
	}
	return datetime.ParseDateFromString(datetime.SiteMonthYear, strings.Split(funcs.StripString(monthDate.FullText()), " ")[0])
}

func (s *Session) GetDate() (datetime.Date, error) {
	root := s.stepToMonth(gorm_models.City{Id: "601", Name: "Saint-Petersburg"}, 0)
	date, err := monthDate(root)
	if err != nil {
		zap.L().Error("Failed to get current date")
		return datetime.Date{}, err
	}
	return date, nil
}

func (s *Session) stepToMonth(city gorm_models.City, delta int) soup.Root {
	var res soup.Root
	if delta == 0 {
		res = s.GetParsedSoup(funcs.Linkify("calendar.do?consularPost=", city.Id))
	} else {
		res = s.GetParsedSoup(funcs.Linkify("calendar.do?month=", strconv.Itoa(delta), "&consularPost=", city.Id))
	}
	s.Date.MoveMonth(delta)
	return res
}

func (s *Session) moveToMonth(city gorm_models.City, date datetime.Date) []soup.Root {
	zap.L().Info("Starting move to month: " + date.Format(datetime.MonthAndYear) + " in city: " + city.Format())
	monthMoveCount := date.Month() - s.Date.Month() + (date.Year()-s.Date.Year())*12
	var responses []soup.Root
	delta := 0
	if monthMoveCount != 0 {
		delta = int(math.Copysign(1, float64(monthMoveCount)))
		for ; monthMoveCount != 0; monthMoveCount -= delta {
			responses = append(responses, s.stepToMonth(city, delta))
			zap.L().Info("Moved to: " + s.Date.Format(datetime.MonthAndYear))
		}
	} else {
		responses = []soup.Root{s.stepToMonth(city, 0)}
	}
	zap.L().Info("Successfully moved to " + date.Format(datetime.MonthAndYear) + " in city: " + city.Format())
	return responses
}

func (s *Session) GetMonthSoup(city gorm_models.City, date datetime.Date) soup.Root {
	zap.L().Info("Starting getting month")
	responses := s.moveToMonth(city, date)
	return responses[len(responses)-1]
}
