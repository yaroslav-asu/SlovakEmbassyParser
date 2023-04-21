package parser

import (
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	gorm_models "main/models/gorm"
	"main/models/gorm/datetime"
	"strings"
)

func (p *Parser) ParseCitiesWithWorkingEmbassies() {
	zap.L().Info("Getting all cities with embassies")
	funcs.Sleep()
	doc := p.Session.GetParsedSoup(funcs.Linkify("consularPost.do"))
	zap.L().Info("Found embassies: ")
	for _, el := range doc.FindAll("option") {
		zap.L().Info(el.Text())
	}
	for _, el := range doc.FindAll("option") {
		city := gorm_models.City{
			Id:   el.Attrs()["value"],
			Name: el.Text(),
		}
		if strings.ToLower(city.Name) == "test" || city.Id == "" {
			continue
		}
		city.StartWorking, city.EndWorking = p.GetEmbassyWorkingMonths(city)
		if city.StartWorking != datetime.NewBlankDate() {
			zap.L().Info("Saving city: " + city.Name + " to db")
			p.SaveToDB(city)
		} else {
			zap.L().Info("Deleting city: " + city.Name + " from db")
			p.DeleteFromDB(city)
		}
	}
	zap.L().Info("Successfully got all cities with embassies")
}
func (p *Parser) isEmbassyWorksInMonth(city gorm_models.City, date datetime.Date) string {
	zap.L().Info("Checking does: " + city.Name + " with id: " + city.Id + " works in: " + date.Format(datetime.MonthAndYear))
	funcs.SleepTime(15, 20)
	doc := p.GetMonthSoup(city, date)
	dayCells := doc.FindAll("td", "class", "calendarMonthCell")
	if len(dayCells) == 0 {
		zap.L().Info("Embassy in " + city.Name + " with id: " + city.Id + " at: " + date.Format(datetime.MonthAndYear) + " doesn't work")
		return "no"
	}
	for dayCell := range dayCells {
		if dayCells[dayCell].Find("strong").Error == nil {
			return "yes"
		}
	}
	return "currently no"
}

func (p *Parser) GetEmbassyWorkingMonths(city gorm_models.City) (datetime.Date, datetime.Date) {
	zap.L().Info("Started checking embassy in " + city.Name + " with id: " + city.Id)
	funcs.SleepTime(15, 20)
	now := datetime.NewDateNow()
	checkingDate := datetime.NewDateYM(now.Year(), now.Month())
	start := datetime.NewBlankDate()
	var end datetime.Date
	for {
		isWorking := p.isEmbassyWorksInMonth(city, checkingDate)
		if isWorking == "yes" {
			if start == datetime.NewBlankDate() {
				start = checkingDate
			}
			end = checkingDate
		}
		if isWorking == "no" {
			return start, end
		}
		checkingDate.MoveMonth(1)
	}
}

func (p *Parser) CheckEmbassyWork(city gorm_models.City) string {
	zap.L().Info("Started checking embassy in " + city.Name + " with id: " + city.Id)
	funcs.Sleep()
	doc := p.Session.GetParsedSoup(funcs.Linkify("calendar.do?consularPost=", city.Id))
	monthCell := doc.Find("td", "class", "calendarMonthCell")
	if monthCell.Error != nil {
		zap.L().Info("Embassy in " + city.Name + " with id: " + city.Id + " doesn't work")
		return "no"
	}
	zap.L().Info("Embassy in " + city.Name + " with id: " + city.Id + " works")
	return "yes"
}

func (p *Parser) CitiesWithWorkingEmbassy() []gorm_models.City {
	var workingCities []gorm_models.City
	p.DB.Find(&workingCities)
	return workingCities
}

func (p *Parser) CityWithWorkingEmbassy(index int) gorm_models.City {
	return p.CitiesWithWorkingEmbassy()[index]
}
