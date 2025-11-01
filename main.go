package main

import (
	"bsnack/config"
	"bsnack/db"
	_ "github.com/joho/godotenv/autoload"
	"log"
)

func main() {
	gormDB, sqlDB, redisClient, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = sqlDB.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err = config.Router(gormDB, redisClient); err != nil {
		log.Fatal(err)
	}
}
