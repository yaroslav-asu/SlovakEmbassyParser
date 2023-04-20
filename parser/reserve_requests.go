package parser

import (
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	gorm_models "main/models/gorm"
	"main/models/gorm/datetime"
	"main/parser/user"
)

func (p *Parser) RunCheckingReserveRequests() {
	for {
		zap.L().Info("Updating reserve requests")
		var reserveRequests []gorm_models.ReserveRequest
		p.DB.DB.Find(&reserveRequests)
		p.DB.DB.Preload("User").Preload("City").Find(&reserveRequests)
		for _, request := range reserveRequests {
			zap.L().Info("Checking reserve request with owner: " + request.User.UserName + " in embassy: " + request.City.Name + " form: " + request.Start.Format(datetime.DateTime) + " to: " + request.End.Format(datetime.DateTime))
			if !request.User.IsReserved {
				for currentDate := request.Start; currentDate.Year() <= request.End.Year() && currentDate.Month() <= request.End.Month(); currentDate.MoveMonth(1) {
					zap.L().Info("Checking date: " + currentDate.Format(datetime.DateOnly))
					for _, dayCell := range p.GetWorkingDaysInMonth(request.City, currentDate) {
						p.DB.DB.Where("id = ?", dayCell.CityId).First(&dayCell.City)
						if dayCell.AvailableReservations > 0 {
							zap.L().Info("Found free reservation at: " + dayCell.City.Name + ", trying to reserve")
							availableReservations, _ := p.GetReservations(dayCell.City, dayCell.Date)
							for _, reservation := range availableReservations {
								userModel := request.User
								reserveRequestOwner := user.NewUser(userModel.UserName, userModel.Password)
								reserveRequestOwner.ReserveDatetime(reservation.City, reservation.Date)
								if reserveRequestOwner.IsReserved() {
									zap.L().Info("User: " + userModel.UserName + "successfully reserved: " + reservation.Date.Format(datetime.DateTime))
								} else {
									zap.L().Error("Failed to reserve user: " + userModel.UserName + " to: " + reservation.Date.Format(datetime.DateTime))
								}
							}
						}
					}
				}
			} else {
				zap.L().Info("User already reserved")
			}
		}
		funcs.SleepTime(10, 90)
	}
}
