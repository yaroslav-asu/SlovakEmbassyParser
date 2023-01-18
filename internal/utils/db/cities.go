package db

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/models"
	"main/parser"
)

func UpdateCities(db *gorm.DB, p parser.Parser) {
	zap.L().Info("Started updating topicalCities with embassies in db")
	topicalCities := p.GetEmbassyCities()
	zap.L().Info("Successfully got topicalCities with embassies")
	for _, city := range topicalCities {
		zap.L().Info("Trying to find or creating city with name: " + city.Name + " and id: " + city.Id + " in db")
		cityCopy := city
		// TODO fix error: when city soft deleted from db, gorm tries to create new one and getting error
		record := db.FirstOrCreate(&cityCopy)
		if record.RowsAffected == 0 {
			zap.L().Info("City with name:" + city.Name + " and id: " + city.Id + " in db doesn't match with current, updating")
			record.Save(&city)
		}
		zap.L().Info("City with name:" + city.Name + " and id: " + city.Id + " up to date")
	}
	DeleteOutdatedCities(db, topicalCities)
}

func DeleteOutdatedCities(db *gorm.DB, topicalCities []models.City) {
	zap.L().Info("Starting to delete outdated cities")
	topicalCitiesMap := make(map[string]bool)
	for _, city := range topicalCities {
		topicalCitiesMap[city.Id] = false
	}
	var dbCities []models.City
	db.Find(&dbCities)
	for _, city := range dbCities {
		_, found := topicalCitiesMap[city.Id]
		if !found {
			zap.L().Info(city.Name + " with id: " + city.Id + "no longer contain an embassy, deleting from db")
			db.Delete(&city)
		}
	}
}
