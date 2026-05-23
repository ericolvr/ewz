package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/ericolvr/ewz/config"
	infraDB "github.com/ericolvr/ewz/internal/infrastructure/database"
	infraHTTP "github.com/ericolvr/ewz/internal/infrastructure/http"
	"github.com/ericolvr/ewz/internal/infrastructure/pipefy"
	"github.com/ericolvr/ewz/internal/interfaces/api"
	"github.com/ericolvr/ewz/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// JSON handler — compatível com Datadog, New Relic, ELK e similares
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

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
	webhookRepo := infraDB.NewWebhookEventRepository(db)
	pipefyClient := pipefy.NewClient()

	clientService := service.NewClientService(clientRepo, pipefyClient)
	webhookService := service.NewWebhookService(clientRepo, webhookRepo, pipefyClient)

	clientHandler := api.NewClientHandler(clientService)
	webhookHandler := api.NewWebhookHandler(webhookService)

	router := gin.Default()
	infraHTTP.SetupRoutes(router, clientHandler, webhookHandler)

	slog.Info("server started", "port", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
