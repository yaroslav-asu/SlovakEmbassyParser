package parser

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"main/models"
	"strings"
)

func (p *Parser) ParseCitiesWithWorkingEmbassies() {
	zap.L().Info("Getting all cities with embassies")
	funcs.RandomSleep()
	res, err := p.getSoup(funcs.Linkify("consularPost.do"))
	if err != nil {
		zap.L().Warn("Failed to connect to /consularPost.do page to get available cities")
	}
	doc := soup.HTMLParse(res)
	for _, el := range doc.FindAll("option") {
		fmt.Println(el.Text())
	}
	for _, el := range doc.FindAll("option") {
		city := models.City{
			Id:   el.Attrs()["value"],
			Name: el.Text(),
		}
		if strings.ToLower(city.Name) == "test" || city.Id == "" {
			continue
		}
		city.StartWorking, city.EndWorking = p.GetEmbassyWorkingMonths(city)
		if city.StartWorking != models.NewBlankDate() {
			p.SaveToDb(city)
		} else {
			p.DeleteFromDb(city)
		}

	}
	zap.L().Info("Successfully got all cities with embassies")
}
func (p *Parser) isEmbassyWorksInMonth(city models.City, date models.Date) string {
	zap.L().Info("Checking does: " + city.Name + " with id: " + city.Id + " work in: " + date.Format(models.MonthAndYear))
	funcs.LongRandomSleep()
	doc := p.GetMonthSoup(city, date)
	dayCells := doc.FindAll("td", "class", "calendarMonthCell")
	if len(dayCells) == 0 {
		zap.L().Info("Embassy in " + city.Name + " with id: " + city.Id + " at: " + date.Format(models.MonthAndYear) + " doesn't work")
		return "no"
	}
	for dayCell := range dayCells {
		if dayCells[dayCell].Find("strong").Error == nil {
			return "yes"
		}
	}
	return "currently no"
}

func (p *Parser) GetEmbassyWorkingMonths(city models.City) (models.Date, models.Date) {
	zap.L().Info("Started checking embassy in " + city.Name + " with id: " + city.Id)
	funcs.LongRandomSleep()
	now := models.NewDateNow()
	checkingDate := models.NewDateYM(now.Year(), now.Month())
	start := models.NewBlankDate()
	var end models.Date
	for {
		isWorking := p.isEmbassyWorksInMonth(city, checkingDate)
		if isWorking == "yes" {
			if start == models.NewBlankDate() {
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

func (p *Parser) CheckEmbassyWork(city models.City) string {
	zap.L().Info("Started checking embassy in " + city.Name + " with id: " + city.Id)
	funcs.RandomSleep()
	doc := p.getParsedSoup(funcs.Linkify("calendar.do?consularPost=", city.Id))
	monthCell := doc.Find("td", "class", "calendarMonthCell")
	if monthCell.Error != nil {
		zap.L().Info("Embassy in " + city.Name + " with id: " + city.Id + " doesn't work")
		return "no"
	}
	zap.L().Info("Embassy in " + city.Name + " with id: " + city.Id + " works")
	return "yes"
}

func (p *Parser) CitiesWithWorkingEmbassy() []models.City {
	var workingCities []models.City
	p.Db.Find(&workingCities)
	return workingCities
}

func (p *Parser) CityWithWorkingEmbassy(index int) models.City {
	return p.CitiesWithWorkingEmbassy()[index]
}
