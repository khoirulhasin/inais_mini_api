package dbs

import (
	"log"

	"os"

	"github.com/jinzhu/gorm"
	"github.com/khoirulhasin/globe_tracker_api/app/models"
)

func OpenPostgres(database string) *gorm.DB {

	databaseDriver := os.Getenv("DATABASE_DRIVER")

	db, err := gorm.Open(databaseDriver, database)
	if err != nil {
		log.Fatalf("%s", err)
	}
	if err := Automigrate(db); err != nil {
		panic(err)
	}
	return db
}

func Automigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Question{}, &models.QuestionOption{}, &models.Answer{}).Error
}
