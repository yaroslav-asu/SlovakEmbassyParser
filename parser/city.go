package parser

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/variables"
	"strings"
)

type City struct {
	title string
	id    string
}

func (p Parser) GetEmbassyCities() []City {
	zap.L().Info("Getting cities where embassies are working")
	p.RandomSleep()
	res, err := soup.GetWithClient(variables.SiteUrl+"consularPost.do", p.client)
	if err != nil {
		zap.L().Warn("Failed to connect to /consularPost.do page to get available cities")
	}
	doc := soup.HTMLParse(res)
	var cities []City
	for _, el := range doc.FindAll("option") {
		city := City{
			el.Text(),
			el.Attrs()["value"],
		}
		if strings.ToLower(city.title) != "test" && city.id != "" {
			if p.CheckEmbassyWork(city) {
				cities = append(cities, city)
			}
		}
	}
	zap.L().Info("Successfully got cities where embassies are working")
	return cities
}

func (p Parser) CheckEmbassyWork(city City) bool {
	zap.L().Info("Started checking embassy in " + city.title + " with id: " + city.id)
	p.RandomSleep()
	res, err := soup.GetWithClient(variables.SiteUrl+"calendar.do?consularPost="+city.id, p.client)
	if err != nil {
		zap.L().Warn("Can't get embassy page of " + city.title + " with id: " + city.id)
	}
	doc := soup.HTMLParse(res)
	monthCell := doc.Find("td", "class", "calendarMonthCell")
	if monthCell.Error != nil {
		zap.L().Info("Embassy in " + city.title + " with id: " + city.id + " doesn't work")
		return false
	}
	zap.L().Info("Embassy in " + city.title + " with id: " + city.id + " works")
	return true
}

func (p Parser) UpdateEmbassyCities() {
	//TODO write to db cities with embassies
}
