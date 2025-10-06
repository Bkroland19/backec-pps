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
	// Drop constraints and recreate optional_vars table to remove primary key
	// This allows duplicate keys with different details
	//DB.Exec("DROP TABLE IF EXISTS optional_vars CASCADE")
	DB.Exec(`CREATE TABLE optional_vars (
		key VARCHAR(255),
		prescriber_type VARCHAR(255),
		intravenous_type VARCHAR(255),
		oral_switch VARCHAR(255),
		number_missed_doses INTEGER,
		missed_doses_reason VARCHAR(255),
		guidelines_compliance VARCHAR(255),
		treatment_type VARCHAR(255),
		parent_key VARCHAR(255)
	)`)

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
