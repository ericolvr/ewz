package api

import (
	"net/http"

	"github.com/ericolvr/ewz/internal/interfaces/dto"
	"github.com/ericolvr/ewz/internal/service"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	service *service.WebhookService
}

func NewWebhookHandler(service *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{service: service}
}

func (h *WebhookHandler) Process(c *gin.Context) {
	var req dto.WebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Process(c.Request.Context(), req.EventID, req.CardID, req.ClientEmail)
	if err != nil {
		if err.Error() == "evento já processado" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "evento processado com sucesso"})
}
