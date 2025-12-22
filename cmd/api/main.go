package main

import (
	"log"
	"os"

	"github.com/NoahFola/travel_app_backend/internal/api"
	"github.com/NoahFola/travel_app_backend/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 2. Init Database
	dbPool := database.InitDB()
	defer dbPool.Close()

	// 3. Init Router (Wires everything together)
	r := api.NewRouter(dbPool)

	// 4. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
