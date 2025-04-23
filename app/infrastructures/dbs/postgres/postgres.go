package postgres

import (
	"fmt"
	"log"

	"os"

	"github.com/jinzhu/gorm"
	"github.com/khoirulhasin/globe_tracker_api/app/models"
)

func Connect() *gorm.DB {
	var (
		databaseUser     = os.Getenv("DATABASE_USER")
		databaseName     = os.Getenv("DATABASE_NAME")
		databaseHost     = os.Getenv("DATABASE_HOST")
		databasePort     = os.Getenv("DATABASE_PORT")
		databasePassword = os.Getenv("DATABASE_PASSWORD")
	)

	var connPostgres = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, databasePort, databaseUser, databaseName, databasePassword)

	databaseDriver := os.Getenv("DATABASE_DRIVER")

	db, err := gorm.Open(databaseDriver, connPostgres)
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
