package parser

import (
	"fmt"
	"main/internal/utils/db"
	"main/internal/utils/funcs"
	gorm_models "main/models/gorm"
	"main/models/gorm/datetime"
	"main/parser/user"
)

func reverseMonths(months []gorm_models.Month) []gorm_models.Month {
	for i := 0; i < len(months)/2; i++ {
		j := len(months) - i - 1
		months[i], months[j] = months[j], months[i]
	}
	return months
}

func (p *Parser) reserveRequestsMonth() []gorm_models.Month {
	var months []gorm_models.Month
	d := db.Connect()
	defer db.Close(d)
	d.Raw("SELECT Distinct \"month_id\" as \"id\", \"date\" FROM \"reserve_requests\" Join months m on m.id = reserve_requests.month_id order by date").Find(&months)
	return months
}

func (p *Parser) firstReserveMonthUser(city gorm_models.City) user.User {
	var userModel gorm_models.User
	p.DB.Raw("SELECT Distinct \"user_id\", \"user_name\", \"password\",  \"city_id\"  FROM \"reserve_requests\" Join users u on reserve_requests.user_id = u.id where city_id = ? order by user_id", city.Id).First(&userModel)
	return user.NewUserFromModel(userModel)
}

func (p *Parser) ParseReserveRequestsInterval(months []gorm_models.Month) {
	for _, month := range months {
		fmt.Println(month.Format())
		var cities []gorm_models.City
		p.DB.Raw("SELECT DISTINCT * from \"reserve_requests\" join cities c on reserve_requests.city_id = c.id where month_id = ?", month.Id).Find(&cities)
		for _, city := range cities {
			dayCells := p.DayCellsWithReservationsInMonth(city, month.Date)
			var userToReserve user.User
			if len(dayCells) > 0 {
				userToReserve = p.firstReserveMonthUser(city)
			}
			for i := 0; i < len(dayCells); i++ {
				availableReservations, _ := p.ReservationsInDay(city, dayCells[i].Date)
				if i < len(dayCells)-1 && userToReserve.ReserveDatetime(city, availableReservations[0].DateTime) {
					userToReserve = p.firstReserveMonthUser(city)
				}
			}
		}
	}

}

func (p *Parser) StartReserveRequestsParsing() {
	// TODO remove kludge
	p.moveToMonth(gorm_models.City{Id: "601", Name: "Saint Peretrburg"}, datetime.Now())
	for months := p.reserveRequestsMonth(); len(months) > 0; months = p.reserveRequestsMonth() {
		p.ParseReserveRequestsInterval(months)
		funcs.Sleep()
		p.ParseReserveRequestsInterval(reverseMonths(months))
		funcs.Sleep()
	}
}
