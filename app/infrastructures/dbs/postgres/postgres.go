package postgres

import (
	"flag"
	"fmt"
	"log"

	"os"

	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	var (
		databaseUser     = os.Getenv("POSTGRES_USER")
		databaseName     = os.Getenv("POSTGRES_NAME")
		databaseHost     = os.Getenv("POSTGRES_HOST")
		databasePort     = os.Getenv("POSTGRES_PORT")
		databasePassword = os.Getenv("POSTGRES_PASSWORD")
		migrate          = flag.Bool("migrate", false, "Run auto migration")
		seed             = flag.Bool("seed", false, "Run database seeding")
	)

	var connPostgres = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, databasePort, databaseUser, databaseName, databasePassword)

	db, err := gorm.Open(postgres.Open(connPostgres), &gorm.Config{})
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Menjalankan ekstensi uuid-ossp jika belum ada
	err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
	if err != nil {
		log.Fatalf("failed to enable uuid-ossp extension: %v", err)
	}

	if err := Automigrate(db); err != nil {
		panic(err)
	}

	flag.Parse()

	// Auto migrate if requested
	if *migrate {
		if err := Automigrate(db); err != nil {
			panic(err)
		}
	}

	// Seed database if requested
	if *seed {
		seeder := NewSeeder(db)
		if err := seeder.SeedAll(); err != nil {
			panic(err)
		}
	}

	log.Println("Database operations completed successfully!")

	return db
}

func Automigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		models.User{},
		models.Role{},
		models.Users2role{},
		models.Menu{},
		models.Menus2role{},
		models.Profile{},
		models.Device{},
		models.Driver{},
		models.Drive{},
		models.Marker{},
		models.Ship{},
		models.MarkerType{},
		models.Cam{},
	)
}
