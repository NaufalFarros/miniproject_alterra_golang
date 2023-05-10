package database

import (
	"fmt"
	"log"
	"os"

	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func DBConnect() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	fmt.Println("Nama database:", os.Getenv("DB_NAME"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.Logger = logger.Default.LogMode(logger.Info)
	db.AutoMigrate(&models.User{}, &models.Items{}, &models.Category{}, &models.Orders{}, &models.OrderItems{}, &models.Table{})

	Database = DbInstance{Db: db}
	return db, nil
}
func DisconnectDB() error {
	sqlDB, err := Database.Db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		return err
	}

	return nil
}
