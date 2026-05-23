package api

import (
	"net/http"

	"github.com/ericolvr/ewz/internal/domain"
	"github.com/ericolvr/ewz/internal/interfaces/dto"
	"github.com/ericolvr/ewz/internal/service"
	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	service *service.ClientService
}

func NewClientHandler(service *service.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

func (h *ClientHandler) Create(c *gin.Context) {
	var req dto.ClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := &domain.Client{
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		RequestType:   req.RequestType,
		AssetValue:    req.AssetValue,
		Status:        "Aguardando Análise",
	}

	if err := h.service.Create(c.Request.Context(), client); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ClientResponse{
		ID:            client.ID,
		CustomerName:  client.CustomerName,
		CustomerEmail: client.CustomerEmail,
		RequestType:   client.RequestType,
		AssetValue:    client.AssetValue,
		Status:        client.Status,
	})
}
