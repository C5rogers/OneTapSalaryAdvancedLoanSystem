package db

import (
	"fmt"
	"log"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/config"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewClient(config *config.Config) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.Server.DBHost,
		config.Server.DBUser,
		config.Server.DBPassword,
		config.Server.DBName,
		config.Server.DBPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}, &models.Customer{}, &models.Transaction{}); err != nil {
		log.Fatal("migration failed:", err)
	}

	return &Database{db}, nil

}
