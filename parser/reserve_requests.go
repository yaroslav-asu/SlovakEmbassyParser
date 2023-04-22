package parser

import (
	gorm_models "main/models/gorm"
	"main/models/gorm/datetime"
	"main/parser/user"
	"os"
)

func (p *Parser) RunCheckingReserveRequests() {
	var userModel gorm_models.User
	p.DB.Where("id = 1").First(&userModel)
	mainUser := user.NewUserFromModel(userModel)
	for !mainUser.Session.User.IsReserved {
		currentDate := datetime.NewDateYM(2023, 6)
		for ; currentDate.Year() <= 2023 && currentDate.Month() <= 8; currentDate.MoveMonth(1) {
			var city gorm_models.City
			city = gorm_models.City{
				Name: "Moscow",
				Id:   "590",
			}
			for _, dayCell := range p.GetWorkingDaysInMonth(city, currentDate) {
				p.DB.Preload("City").Find(&dayCell)
				if dayCell.AvailableReservations > 0 {
					availableReservations, _ := p.GetReservations(city, dayCell.Date)
					for _, reservation := range availableReservations {
						if mainUser.ReserveDatetime(city, reservation.Date) {
							os.Exit(1)
						}
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
			for _, dayCell := range p.GetWorkingDaysInMonth(city, currentDate) {
				p.DB.Preload("City").Find(&dayCell)
				if dayCell.AvailableReservations > 0 {
					availableReservations, _ := p.GetReservations(city, dayCell.Date)
					for _, reservation := range availableReservations {
						if mainUser.ReserveDatetime(city, reservation.Date) {
							os.Exit(1)
						}
					}
				}
			}
		}
	}
}
