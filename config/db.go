package config

import (
	"arithmetic-calculator/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	log.Println("Connecting to postgres database...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database server: %v", err)
	}

	var exists bool
	err = db.Raw("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = ?)", os.Getenv("DB_NAME")).Scan(&exists).Error
	if err != nil {
		log.Fatalf("Failed to check if database exists: %v", err)
	}

	if !exists {
		log.Println("Database does not exist, creating database...")
		dbtest := db.Exec("CREATE DATABASE " + os.Getenv("DB_NAME"))
		if dbtest == nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		log.Println("Database created successfully!")
	}

	dsnWithDB := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	DB, err = gorm.Open(postgres.Open(dsnWithDB), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	DB.AutoMigrate(&models.User{}, &models.Operation{}, &models.Record{})
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "dev" || appEnv == "prod" {
		SeedDatabase()
	}
	log.Println("Database connection established!")
}

func SeedDatabase() {
	var userCount int64
	var operationCount int64

	if err := DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		log.Fatalf("Error counting users: %v", err)
	}
	if userCount == 0 {
		users := []models.User{
			{
				Username: "testuser1@example.com",
				Password: "$2a$10$KIXW/x75r8yL.IK1ZyQeOO1j7xHYQaQRV8MBICo.e5MTeQUH96C/q",
				Balance:  100.00,
				Status:   "active",
			}, // pass: password123
			{
				Username: "testuser2@example.com",
				Password: "$2a$10$KIXW/x75r8yL.IK1ZyQeOO1j7xHYQaQRV8MBICo.e5MTeQUH96C/q",
				Balance:  150.00,
				Status:   "active",
			}, // pass: password123
		}
		DB.Create(&users)
		log.Println("Seeded users.")
	}

	if err := DB.Model(&models.Operation{}).Count(&operationCount).Error; err != nil {
		log.Fatalf("Error counting operations: %v", err)
	}
	if operationCount == 0 {
		operations := []models.Operation{
			{Type: "addition", Cost: 1.00},
			{Type: "subtraction", Cost: 1.00},
			{Type: "multiplication", Cost: 1.50},
			{Type: "division", Cost: 2.00},
			{Type: "square_root", Cost: 2.50},
			{Type: "random_string", Cost: 3.00},
		}
		DB.Create(&operations)
		log.Println("Seeded operations.")
	}
}
