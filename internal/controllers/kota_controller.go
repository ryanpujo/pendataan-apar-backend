package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/pendataan-apar-backend/internal/models"
	"github.com/ryanpujo/pendataan-apar-backend/internal/services"
)

type Kota interface {
	Create(c *gin.Context)
}

type KotaController struct {
	kotaService services.Kota
}

func NewKotaController(kota services.Kota) *KotaController {
	return &KotaController{
		kotaService: kota,
	}
}

type CreateKotaResponse struct {
	ID      uint   `json:"id,omitempty"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func (kc *KotaController) Create(c *gin.Context) {
	var payload models.CreateKotaRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, CreateKotaResponse{
			Message: "validation error",
			Error:   err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	id, err := kc.kotaService.Create(ctx, &payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, CreateKotaResponse{
			Message: "Failed to create kota",
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, CreateKotaResponse{
		ID: id,
		Message: "Berhasil", 
	})
}
