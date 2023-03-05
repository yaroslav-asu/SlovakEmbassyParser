package parser

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"main/models"
	"strconv"
	"strings"
	"time"
)

func (p *Parser) moveToMonth(city models.City, month int, year int) string {
	zap.L().Info("Starting move to another month")
	res, err := p.GetSoup(funcs.Linkefy("calendar.do?month=", strconv.Itoa(month-p.Date.Month()+(year-p.Date.Year())*12), "&consularPost=", city.Id))
	if err != nil {
		zap.L().Warn("Failed move to month")
		return ""
	}
	p.Date.ChangeMonth(month)
	p.Date.ChangeYear(year)
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
		dateText := funcs.StripString(dateNode.Text()) + strconv.Itoa(p.Date.Year())
		availableReservations := ParseReservationData(reservationData)
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

func ParseReservationData(data string) int {
	for _, s := range []string{"[", "]"} {
		data = strings.Replace(data, s, "", -1)
	}
	parsedNumbers := strings.Split(data, "/")
	reservedNum, err := strconv.Atoi(parsedNumbers[0])
	if err != nil {
		zap.L().Error("Failed to parse count of reserved times")
		return 0
	}
	totalNum, err := strconv.Atoi(parsedNumbers[1])
	if err != nil {
		zap.L().Error("Failed to parse amount of reservations")
		return 0
	}
	return totalNum - reservedNum
}

func (p *Parser) ParseDay(city models.City, date models.Date) models.AvailableReservations {
	var availableReservations models.AvailableReservations
	dateString := date.Format("02.01.2006")
	zap.L().Info("Started parsing available reservations of: " + dateString + " in " + city.Name)
	res, err := p.GetSoup(funcs.Linkefy("calendarDay.do?day=", dateString, "&consularPostId=", city.Id))
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
		date.ChangeHour(availableTime.Hour())
		date.ChangeMinutes(availableTime.Minute())
		availableReservation := models.AvailableReservation{
			CityId: city.Id,
			Date:   date,
		}
		availableReservations = append(availableReservations, availableReservation)
	}
	zap.L().Info("Finished parsing available reservations of: " + dateString + " in " + city.Name)
	return availableReservations
}
