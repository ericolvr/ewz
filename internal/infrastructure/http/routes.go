package http

import (
	"github.com/ericolvr/ewz/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	clientHandler *api.ClientHandler,
	webhookHandler *api.WebhookHandler,
) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/clients", clientHandler.Create)
		v1.GET("/clients/:email", clientHandler.GetByEmail)
		v1.POST("/webhooks/pipefy/card-updated", webhookHandler.Process)
	}
}
