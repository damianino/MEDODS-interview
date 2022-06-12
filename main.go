package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbClient := InitDbConnection()
	defer dbClient.Disconnect(context.Background())

	r := gin.Default()

	r.POST("/login", loginHandler(dbClient))
	r.POST("/token", tokenHandler(dbClient))

	r.Run()
}
