package parser

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/variables"
	"main/models"
	"strings"
)

func (p Parser) GetEmbassyCities() []models.City {
	zap.L().Info("Getting all cities with embassies")
	p.RandomSleep()
	res, err := soup.GetWithClient(variables.SiteUrl+"consularPost.do", p.client)
	if err != nil {
		zap.L().Warn("Failed to connect to /consularPost.do page to get available cities")
	}
	doc := soup.HTMLParse(res)
	var cities []models.City
	for _, el := range doc.FindAll("option") {
		city := models.City{
			Id:   el.Attrs()["value"],
			Name: el.Text(),
		}
		city.Working = p.CheckEmbassyWork(city)
		if strings.ToLower(city.Name) != "test" && city.Id != "" {
			cities = append(cities, city)
		}
	}
	zap.L().Info("Successfully got all cities with embassies")
	return cities
}
func (p Parser) CheckEmbassyWork(city models.City) bool {
	zap.L().Info("Started checking embassy in " + city.Name + " with id: " + city.Id)
	p.RandomSleep()
	res, err := soup.GetWithClient(variables.SiteUrl+"calendar.do?consularPost="+city.Id, p.client)
	if err != nil {
		zap.L().Warn("Can't get embassy page of " + city.Name + " with id: " + city.Id)
	}
	doc := soup.HTMLParse(res)
	monthCell := doc.Find("td", "class", "calendarMonthCell")
	if monthCell.Error != nil {
		zap.L().Info("Embassy in " + city.Name + " with id: " + city.Id + " doesn't work")
		return false
	}
	zap.L().Info("Embassy in " + city.Name + " with id: " + city.Id + " works")
	return true
}
