package parser

import (
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	gorm_models "main/models/gorm"
	"main/models/gorm/datetime"
	"strconv"
	"strings"
	"time"
)

func availableReservationsInDay(data string) int {
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

func (p *Parser) GetReservations(city gorm_models.City, date datetime.Date) ([]gorm_models.Reservation, []gorm_models.Reservation) {
	funcs.Sleep()
	var availableReservations, unavailableReservations []gorm_models.Reservation
	dateString := date.Format(datetime.DateOnly)
	zap.L().Info("Started parsing available reservations of: " + dateString + " in " + city.Name)
	doc := p.Session.GetParsedSoup(funcs.Linkify("calendarDay.do?day=", dateString, "&consularPostId=", city.Id))
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
			zap.L().Fatal("Error on trying to parse time: " + textTime)
			continue
		}
		date.SetHour(parsedTime.Hour())
		date.ChangeMinutes(parsedTime.Minute())
		reservation := gorm_models.Reservation{
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
