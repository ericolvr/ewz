package main

import (
	"database/sql"
	"log"

	"github.com/ericolvr/ewz/config"
	infraDB "github.com/ericolvr/ewz/internal/infrastructure/database"
	infraHTTP "github.com/ericolvr/ewz/internal/infrastructure/http"
	"github.com/ericolvr/ewz/internal/interfaces/api"
	"github.com/ericolvr/ewz/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	clientRepo := infraDB.NewClientRepository(db)
	clientService := service.NewClientService(clientRepo)
	clientHandler := api.NewClientHandler(clientService)

	router := gin.Default()
	infraHTTP.SetupRoutes(router, clientHandler)

	log.Printf("server running on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
