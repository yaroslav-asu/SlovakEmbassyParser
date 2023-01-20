package parser

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"main/models"
	"strconv"
)

func (p *Parser) moveToMonth(city models.City, month int, year int) string {
	zap.L().Info("Starting move to another month")
	a := strconv.Itoa(month - p.Month + (year-p.Year)*12)
	fmt.Println(a)
	link := funcs.Linkefy("calendar.do?month=", strconv.Itoa(month-p.Month+(year-p.Year)*12), "&consularPost=", city.Id)
	res, err := p.Get(link)
	if err != nil {
		zap.L().Warn("Failed move to month")
		return ""
	}
	p.Month = month
	p.Year = year
	zap.L().Info("Successfully moved to month")
	return res
}

func (p *Parser) GetMonth(city models.City, month int, year int) string {
	zap.L().Info("Starting getting month")
	return p.moveToMonth(city, month, year)
}

func (p *Parser) ParseMonth(city models.City, month int, year int) []models.DayCell {
	var dayCells []models.DayCell
	res := p.GetMonth(city, month, year)
	monthCell := soup.HTMLParse(res).FindAll("td", "class", "calendarMonthCell")
	for _, el := range monthCell {
		freeSpaceNode := el.Find("font")
		if freeSpaceNode.Error != nil {
			continue
		}
		reservationData, dateNode := funcs.StripString(freeSpaceNode.Text()), freeSpaceNode.Find("strong")
		if dateNode.Error != nil {
			continue
		}
		dateText := funcs.StripString(dateNode.Text())
		availableReservations := ParseReservationData(reservationData)
		date := ParseMonthCellDate(dateText, p.Year)
		dayCell := models.DayCell{
			AvailableReservations: availableReservations,
			Date:                  date,
		}
		dayCells = append(dayCells, dayCell)
	}
	return dayCells
}
