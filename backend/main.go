package main

import (
	"association/config"
	"association/database"
	"association/routes"
	"log"
)

func main() {
	log.Println("Starting Association Management System...")

	config.InitConfig()

	database.InitDatabase()

	r := routes.SetupRouter()

	port := config.AppConfig.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
