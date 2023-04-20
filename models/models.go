package models

import (
	"main/internal/utils/db"
)

type DbModel interface {
	SaveToDB(db *db.DB)
	DeleteFromDB(db *db.DB)
}
