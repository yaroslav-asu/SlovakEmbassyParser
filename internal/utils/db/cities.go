package db

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/models"
	"main/parser"
)

func UpdateCities(db *gorm.DB, p parser.Parser) {
	zap.L().Info("Started updating cities with embassies in db")
	cities := p.GetEmbassyCities()
	zap.L().Info("Successfully got cities with embassies")
	for _, city := range cities {
		zap.L().Info("Trying to find city with name: " + city.Name + " and id: " + city.Id + " in db")
		result := models.City{}
		record := db.Model(&models.City{Id: city.Id}).First(&result)
		if record.Error != nil {
			if errors.Is(record.Error, gorm.ErrRecordNotFound) {
				zap.L().Info("City with name: " + city.Name + " and id: " + city.Id + " not found in db")
				db.Create(&city)
				continue
			} else {
				zap.L().Error("Got unknown error: " + record.Error.Error())
			}
		}
		if result.Name != city.Name || result.Working != city.Working {
			record.Save(&city)
			zap.L().Info("City with name: " + city.Name + " and id: " + city.Id + " successfully updated")
		}
		zap.L().Info("City with name: " + city.Name + " and id: " + city.Id + " not needed in update")
	}
}
