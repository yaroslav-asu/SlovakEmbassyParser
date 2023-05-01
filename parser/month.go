package parser

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	gorm_models "main/models/gorm"
	"main/models/gorm/datetime"
	"math"
	"strconv"
)

//TODO complit functions
//func (p *Parser) monthDate(root soup.Root) (datetime.Date, error) {
//	monthDate := root.Find("td", "class", "calendarMonthLabel")
//	if monthDate.Error != nil {
//		return datetime.Date{}, monthDate.Error
//	}
//	datetime.ParseDateFromString() strings.Split(monthDate.FullText(), " ")[0]
//	return
//}
//
//func (p *Parser) CurrentMonth() {
//	return p.stepToMonth(gorm_models.City{Id: "601", Name: "Saint-Petersburg"}, 0)
//}

func (p *Parser) stepToMonth(city gorm_models.City, delta int) soup.Root {
	var res soup.Root
	if delta == 0 {
		res = p.Session.GetParsedSoup(funcs.Linkify("calendar.do?consularPost=", city.Id))
	} else {
		res = p.Session.GetParsedSoup(funcs.Linkify("calendar.do?month=", strconv.Itoa(delta), "&consularPost=", city.Id))
	}
	p.Date.MoveMonth(delta)
	return res
}

func (p *Parser) moveToMonth(city gorm_models.City, date datetime.Date) []soup.Root {
	zap.L().Info("Starting move to month: " + date.Format(datetime.MonthAndYear) + " in city: " + city.Format())
	monthMoveCount := date.Month() - p.Date.Month() + (date.Year()-p.Date.Year())*12
	var responses []soup.Root
	delta := 0
	if monthMoveCount != 0 {
		delta = int(math.Copysign(1, float64(monthMoveCount)))
		for ; monthMoveCount != 0; monthMoveCount -= delta {
			responses = append(responses, p.stepToMonth(city, delta))
		}
	} else {
		responses = []soup.Root{p.stepToMonth(city, 0)}
	}
	zap.L().Info("Successfully moved to " + date.Format(datetime.MonthAndYear) + " in city: " + city.Format())
	return responses
}

func (p *Parser) getMonthSoup(city gorm_models.City, date datetime.Date) soup.Root {
	zap.L().Info("Starting getting month")
	responses := p.moveToMonth(city, date)
	return responses[len(responses)-1]
}

func (p *Parser) getWorkingDaysInMonth(city gorm_models.City, date datetime.Date) []gorm_models.DayCell {
	zap.L().Info("Started to get day cells")
	funcs.Sleep()
	var dayCells []gorm_models.DayCell
	res := p.getMonthSoup(city, date)
	monthCell := res.FindAll("td", "class", "calendarMonthCell")
	for _, el := range monthCell {
		freeSpaceNode := el.Find("font")
		if freeSpaceNode.Error != nil {
			continue
		}
		reservationData, dateNode := funcs.StripString(freeSpaceNode.Text()), freeSpaceNode.Find("strong")
		if dateNode.Error != nil {
			continue
		}
		dateText := funcs.StripString(dateNode.Text()) + strconv.Itoa(date.Year())
		availableReservations := availableReservationsInDay(reservationData)
		date := datetime.ParseDateFromString(dateText)
		dayCell := gorm_models.DayCell{
			AvailableReservations: availableReservations,
			CityId:                city.Id,
			Date:                  date,
		}
		dayCells = append(dayCells, dayCell)
	}
	zap.L().Info("Finished to get day cells")
	return dayCells
}

func (p *Parser) DayCellsWithReservationsInMonth(city gorm_models.City, date datetime.Date) []gorm_models.DayCell {
	dayCells := p.getWorkingDaysInMonth(city, date)
	var DayCellsWithReservations []gorm_models.DayCell
	for _, dayCell := range dayCells {
		if dayCell.AvailableReservations > 0 {
			DayCellsWithReservations = append(DayCellsWithReservations, dayCell)
		}
	}
	return DayCellsWithReservations
}
