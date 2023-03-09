package parser

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"main/models"
	"strconv"
)

func (p *Parser) moveToMonth(city models.City, date models.Date) soup.Root {
	zap.L().Info("Starting move to month: " + date.Format(models.MonthAndYear))
	res := p.getParsedSoup(funcs.Linkify("calendar.do?month=", strconv.Itoa(date.Month()-p.Date.Month()+(date.Year()-p.Date.Year())*12), "&consularPost=", city.Id))
	p.Date.SetMonth(date.Month())
	p.Date.SetYear(date.Year())
	zap.L().Info("Successfully moved to " + date.Format(models.MonthAndYear))
	return res
}

func (p *Parser) GetMonthSoup(city models.City, date models.Date) soup.Root {
	zap.L().Info("Starting getting month")
	return p.moveToMonth(city, date)
}

func (p *Parser) GetWorkingDaysInMonth(city models.City, date models.Date) []models.DayCell {
	zap.L().Info("Started to get day cells")
	funcs.Sleep()
	var dayCells []models.DayCell
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
		date := models.ParseDateFromString(dateText)
		dayCell := models.DayCell{
			AvailableReservations: availableReservations,
			CityId:                city.Id,
			Date:                  date,
		}
		dayCells = append(dayCells, dayCell)
	}
	zap.L().Info("Finished to get day cells")
	return dayCells
}
