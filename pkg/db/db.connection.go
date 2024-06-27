package db

import (
	"fmt"
	"log"
	"os"
	"trade/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {

	DB := Connect()
	db, err := DB.DB()
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	err = db.Ping()
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	log.Println("Successfully connected to the database!")

	err = DB.AutoMigrate(&models.CryptoHistory{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

func Connect() *gorm.DB {
	dsn := "host=" + os.Getenv("POSTGRES_HOST") + " port=" + os.Getenv("POSTGRES_PORT") + " user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=" + os.Getenv("POSTGRES_SSL_MODE") + ""

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("dns : ", dsn)
		log.Printf("failed to connect to database: %v", err)
	}

	return DB
}
