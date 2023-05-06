package parser

import (
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	gorm_models "main/models/gorm"
	"main/models/gorm/datetime"
	"strconv"
)

func (p *Parser) monthWorkDays(city gorm_models.City, date datetime.Date) []gorm_models.DayCell {
	zap.L().Info("Started to get day cells")
	funcs.Sleep()
	var dayCells []gorm_models.DayCell
	res := p.Session.GetMonthSoup(city, date)
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
		date, err := datetime.ParseDateFromString(datetime.DateOnly, dateText)
		if err != nil {
			zap.L().Error("Failed to get parsed date, continuing")
			continue
		}
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
	dayCells := p.monthWorkDays(city, date)
	var DayCellsWithReservations []gorm_models.DayCell
	for _, dayCell := range dayCells {
		if dayCell.AvailableReservations > 0 {
			DayCellsWithReservations = append(DayCellsWithReservations, dayCell)
		}
	}
	return DayCellsWithReservations
}
