package parser

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/datetime"
	"main/internal/utils/funcs"
	gorm_models "main/models/gorm"
	"strconv"
)

func (p *Parser) moveToMonth(city gorm_models.City, date datetime.Date) soup.Root {
	zap.L().Info("Starting move to month: " + date.Format(datetime.MonthAndYear))
	res := p.Session.GetParsedSoup(funcs.Linkify("calendar.do?month=", strconv.Itoa(date.Month()-p.Date.Month()+(date.Year()-p.Date.Year())*12), "&consularPost=", city.Id))
	p.Date.SetMonth(date.Month())
	p.Date.SetYear(date.Year())
	zap.L().Info("Successfully moved to " + date.Format(datetime.MonthAndYear))
	return res
}

func (p *Parser) GetMonthSoup(city gorm_models.City, date datetime.Date) soup.Root {
	zap.L().Info("Starting getting month")
	return p.moveToMonth(city, date)
}

func (p *Parser) GetWorkingDaysInMonth(city gorm_models.City, date datetime.Date) []gorm_models.DayCell {
	zap.L().Info("Started to get day cells")
	funcs.Sleep()
	var dayCells []gorm_models.DayCell
	res := p.GetMonthSoup(city, date)
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
		availableReservations := AvailableReservationsInDay(reservationData)
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
