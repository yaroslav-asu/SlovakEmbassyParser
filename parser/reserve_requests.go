package parser

import (
	"fmt"
	gorm_models "main/models/gorm"
	"main/models/gorm/datetime"
	"main/parser/user"
)

func (p *Parser) RunCheckingReserveRequests() {
	var userModel gorm_models.User
	p.DB.Where("id = 1").First(&userModel)
	mainUser := user.NewUserFromModel(userModel)
	for {
		currentDate := datetime.NewDateYM(2023, 6)
		for ; currentDate.Year() <= 2023 && currentDate.Month() <= 7; currentDate.MoveMonth(1) {
			var city gorm_models.City
			city = gorm_models.City{
				Name: "Moscow",
				Id:   "590",
			}
			fmt.Println(currentDate.Format(datetime.DateTime), city.Name)
			for _, dayCell := range p.GetWorkingDaysInMonth(city, currentDate) {
				if dayCell.AvailableReservations > 0 {
					availableReservations, _ := p.GetReservations(city, currentDate)
					for _, reservation := range availableReservations {
						mainUser.ReserveDatetime(city, reservation.Date)
					}
				}
			}
		}
		currentDate.MoveMonth(-1)
		for ; currentDate.Year() >= 2023 && currentDate.Month() >= 6; currentDate.MoveMonth(-1) {
			var city gorm_models.City
			city = gorm_models.City{
				Name: "Saint Petersburg",
				Id:   "601",
			}
			fmt.Println(currentDate.Format(datetime.DateTime), city.Name)

			for _, dayCell := range p.GetWorkingDaysInMonth(city, currentDate) {
				if dayCell.AvailableReservations > 0 {
					availableReservations, _ := p.GetReservations(city, currentDate)
					for _, reservation := range availableReservations {
						mainUser.ReserveDatetime(city, reservation.Date)
					}
				}
			}
		}
	}
}
