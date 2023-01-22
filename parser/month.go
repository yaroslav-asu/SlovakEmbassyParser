package parser

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"main/models"
	"strconv"
	"time"
)

func (p *Parser) moveToMonth(city models.City, month int, year int) string {
	zap.L().Info("Starting move to another month")
	res, err := p.Get(funcs.Linkefy("calendar.do?month=", strconv.Itoa(month-p.Month+(year-p.Year)*12), "&consularPost=", city.Id))
	if err != nil {
		zap.L().Warn("Failed move to month")
		return ""
	}
	p.Month = month
	p.Year = year
	zap.L().Info("Successfully moved to month")
	return res
}

func (p *Parser) GetMonthSoup(city models.City, month int, year int) string {
	zap.L().Info("Starting getting month")
	return p.moveToMonth(city, month, year)
}

func (p *Parser) GetDayCells(city models.City, month int, year int) []models.DayCell {
	zap.L().Info("Started to get day cells")
	var dayCells []models.DayCell
	res := p.GetMonthSoup(city, month, year)
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
			CityId:                city.Id,
			Date:                  date,
		}
		dayCells = append(dayCells, dayCell)
	}
	zap.L().Info("Finished to get day cells")
	return dayCells
}

func (p *Parser) ParseAvailableReservations(city models.City, date time.Time) models.AvailableReservations {
	var availableReservations models.AvailableReservations
	dateString := date.Format("02.01.2006")
	zap.L().Info("Started parsing available reservations of: " + dateString + " in " + city.Name)
	res, err := p.Get(funcs.Linkefy("calendarDay.do?day=", dateString, "&consularPostId=", city.Id))
	if err != nil {
		zap.L().Warn("Got error, from getting date page of: " + dateString + " in " + city.Name + ":\n" + err.Error())
	}
	doc := soup.HTMLParse(res)
	trs := doc.FindAll("tr")
	for _, tr := range trs {
		conditionNode := tr.Find("td", "class", "calendarDayTableRow")
		if conditionNode.Error != nil || funcs.StripString(conditionNode.Text()) == "full" {
			continue
		}
		timeNode := tr.Find("td", "class", "calendarDayTableDateColumn")
		if timeNode.Error != nil {
			continue
		}
		availableTimeText := funcs.StripString(timeNode.Text())
		availableTime, err := time.Parse("15:04", availableTimeText)
		if err != nil {

		}
		fullDate := date.Add(time.Hour*time.Duration(availableTime.Hour()) + time.Minute*time.Duration(availableTime.Minute()))
		availableReservation := models.AvailableReservation{
			CityId: city.Id,
			Date:   fullDate,
		}
		availableReservations.Reservations = append(availableReservations.Reservations, availableReservation)
	}
	zap.L().Info("Finished parsing available reservations of: " + dateString + " in " + city.Name)
	return availableReservations
}
