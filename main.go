package main

import (
	"log"
	"os"
	router "trade/Router"
	"trade/pkg/db"

	"github.com/joho/godotenv"
)

func init() {

	if len(os.Getenv("POSTGRES_HOST")) == 0 {
		err := godotenv.Load()
		if err != nil {
			log.Println("Failed to parse the env file")
			panic(err)
		}

	}
}

func init() {
	db.Init()
}

func main() {
	router.Router()
}
