package main

import (
	"github.com/AdrianOrlow/links-api/app"
	"github.com/AdrianOrlow/links-api/config"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	config := config.LoadConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":" + config.Port)
}