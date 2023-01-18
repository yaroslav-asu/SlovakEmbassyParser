package parser

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"main/models"
	"strings"
)

func (p Parser) GetEmbassyCities() []models.City {
	zap.L().Info("Getting all cities with embassies")
	p.RandomSleep()
	res, err := p.Get(funcs.Linkefy("consularPost.do"))
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
	res, err := p.Get(funcs.Linkefy("calendar.do?consularPost=", city.Id))
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

func (p Parser) UpdateCities() {
	zap.L().Info("Started updating topicalCities with embassies in db")
	topicalCities := p.GetEmbassyCities()
	zap.L().Info("Successfully got topicalCities with embassies")
	for _, city := range topicalCities {
		zap.L().Info("Trying to find or creating city with name: " + city.Name + " and id: " + city.Id + " in db")
		cityCopy := city
		// TODO fix error: when city soft deleted from db, gorm tries to create new one and getting error
		record := p.db.FirstOrCreate(&cityCopy)
		if record.RowsAffected == 0 {
			zap.L().Info("City with name:" + city.Name + " and id: " + city.Id + " in db doesn't match with current, updating")
			record.Save(&city)
		}
		zap.L().Info("City with name:" + city.Name + " and id: " + city.Id + " up to date")
	}
	p.DeleteOutdatedCities(topicalCities)
}

func (p Parser) DeleteOutdatedCities(topicalCities []models.City) {
	zap.L().Info("Starting to delete outdated cities")
	topicalCitiesMap := make(map[string]bool)
	for _, city := range topicalCities {
		topicalCitiesMap[city.Id] = false
	}
	var dbCities []models.City
	p.db.Find(&dbCities)
	for _, city := range dbCities {
		_, found := topicalCitiesMap[city.Id]
		if !found {
			zap.L().Info(city.Name + " with id: " + city.Id + "no longer contain an embassy, deleting from db")
			p.db.Delete(&city)
		}
	}
}
func (p Parser) WorkingCityByIndex(index int) models.City {
	var workingCities []models.City
	p.db.Where("working = ?", true).Find(&workingCities)
	return workingCities[index]
}
