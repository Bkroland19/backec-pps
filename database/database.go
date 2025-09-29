package database

import (
	"log"
	"point-prevalence-survey/config"
	"point-prevalence-survey/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.LoadConfig()

	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("Database connected successfully")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.Patient{},
		&models.Antibiotic{},
		&models.AntibioticDetails{},
		&models.Indication{},
		&models.OptionalVar{},
		&models.Specimen{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed successfully")
}

func GetDB() *gorm.DB {
	return DB
}
