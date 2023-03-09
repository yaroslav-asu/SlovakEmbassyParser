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

func AvailableReservationsInDay(data string) int {
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

func (p *Parser) GetReservations(city models.City, date models.Date) (models.Reservations, models.Reservations) {
	funcs.Sleep()
	var availableReservations, unavailableReservations models.Reservations
	dateString := date.Format(models.BasicDate)
	zap.L().Info("Started parsing available reservations of: " + dateString + " in " + city.Name)
	res, err := p.getSoup(funcs.Linkify("calendarDay.do?day=", dateString, "&consularPostId=", city.Id))
	if err != nil {
		zap.L().Warn("Got error, from getting date page of: " + dateString + " in " + city.Name + ":\n" + err.Error())
	}
	doc := soup.HTMLParse(res)
	trs := doc.FindAll("tr")
	for _, tr := range trs {
		conditionNode := tr.Find("td", "class", "calendarDayTableRow")
		if conditionNode.Error != nil {
			continue
		}
		timeNode := tr.Find("td", "class", "calendarDayTableDateColumn")
		if timeNode.Error != nil {
			continue
		}
		textTime := funcs.StripString(timeNode.Text())
		parsedTime, err := time.Parse("15:04", textTime)
		if err != nil {
			zap.L().Error("Error on trying to parse time")
			continue
		}
		date.SetHour(parsedTime.Hour())
		date.ChangeMinutes(parsedTime.Minute())
		reservation := models.Reservation{
			CityId: city.Id,
			Date:   date,
		}
		if funcs.StripString(conditionNode.Text()) == "full" {
			unavailableReservations = append(unavailableReservations, reservation)
		} else {
			availableReservations = append(availableReservations, reservation)
		}
	}
	zap.L().Info("Finished parsing available reservations of: " + dateString + " in " + city.Name)
	return availableReservations, unavailableReservations
}

func (p *Parser) ParseMonthReservations(city models.City, date models.Date) {
	workingDays := p.GetWorkingDaysInMonth(city, date)
	for i := range workingDays {
		availableReservations, unavailableReservations := p.GetReservations(city, workingDays[i].Date)
		if workingDays[i].AvailableReservations > 0 {
			p.SaveToDB(availableReservations)
		}
		p.DeleteFromDB(unavailableReservations)
		funcs.SleepTime(5, 10)
	}
}

func (p *Parser) ParseReservations() {
	cities := p.CitiesWithWorkingEmbassy()
	for city := range cities {
		currentCity := cities[city]
		for date := currentCity.StartWorking; date != currentCity.EndWorking; date.MoveMonth(1) {
			p.ParseMonthReservations(currentCity, date)
			funcs.SleepTime(15, 30)
		}
		funcs.SleepTime(30, 60)
	}
}
