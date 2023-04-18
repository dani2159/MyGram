package main

import (
	"MyGramAPI/app/router"
	"MyGramAPI/pkg/database"
	"log"
	"os"

	"github.com/joho/godotenv"
)


func envinit() {
	if os.Getenv("ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func main() {
	envinit()

	// db.Connect()
	database.GetDB()
	route := router.InitRouter()
	route.Run()

}