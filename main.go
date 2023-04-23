package main

import (
	"MyGramAPI/app/routers"
	"MyGramAPI/pkg/helpers"
	"log"

	"github.com/joho/godotenv"
)

// @title MyGram API
// @version 1.0
// @description This is an API for MyGram APP. To use all of the services, please login first and get the token.
// @description Once you've completed the previous steps, you'll need to locate the "Authorize" button on the right-hand side of the screen and click it. This will trigger a pop-up window to appear, in which you should enter your token preceded by the word "Bearer". For instance, your token might look something like "eyJhbGciOiJIUzI1...", so you would enter "Bearer eyJhbGciOiJIUzI1..." into the designated field.

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  danisetiawan609@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @host localhost:8082
// @BasePath /
// @swagg.NoModels

func init() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
}

func main() {

	helpers.InitCloudinary()
	routers.StartServer().Run()
}
