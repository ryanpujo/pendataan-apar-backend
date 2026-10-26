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
	FindByID(c *gin.Context)
}

type KotaController struct {
	kotaService services.Kota
}

type Response struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func NewKotaController(kota services.Kota) *KotaController {
	return &KotaController{
		kotaService: kota,
	}
}

type CreateKotaResponse struct {
	ID      uint   `json:"id,omitempty"`
	Response
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


type FindKotaByIDResponse struct {
	Kota models.Kota `json:"kota"`
	Response
}
func (kc *KotaController) FindByID(c *gin.Context) {
	var uri struct{
		ID uint `uri:"id" binding:"required,gt=0"`
	}

	if err := c.ShouldBindUri(&uri);err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, FindKotaByIDResponse{
			Message: "validation error",
			Error: err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	kota, err := kc.kotaService.FindByID(ctx, uri.ID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, FindKotaByIDResponse{
			Message: "failed to find kota",
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, FindKotaByIDResponse{
		Kota: *kota,
		Message: "Berhasil",
	})
}