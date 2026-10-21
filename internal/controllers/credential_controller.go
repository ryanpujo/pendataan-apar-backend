package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/pendataan-apar-backend/internal/models"
	"github.com/ryanpujo/pendataan-apar-backend/internal/services"
)

type Credential interface {
	Register(c *gin.Context)
	FindByID(c *gin.Context)
}

type CredentialController struct {
	credentialService services.Credential
}

func NewCredentialController(credential services.Credential) *CredentialController {
	return &CredentialController{
		credentialService: credential,
	}
}

type RegisterResponse struct {
	ID      uint   `json:"id,omitempty"`
	Message string `json:"message"`
	Error   error  `json:"error,omitempty"`
}

func (cc *CredentialController) Register(c *gin.Context) {
	var payload models.RegisterCredentialRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, RegisterResponse{
			Message: "validation error",
			Error:   err,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	id, err := cc.credentialService.Register(ctx, &payload)
	if err != nil {
		if err == context.DeadlineExceeded {
			c.AbortWithStatusJSON(http.StatusBadRequest, RegisterResponse{
				Message: "request timeout",
				Error:   err,
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, RegisterResponse{
			Message: "failed to create credential",
			Error:   err,
		})
		return
	}

	json := RegisterResponse{
		ID:      id,
		Message: "Berhasil",
	}
	c.JSON(http.StatusCreated, json)
}

type FindByIDResponse struct {
	Credential *models.Credential `json:"credential,omitempty"`
	Message string `json:"message"`
	Error string `json:"error,omitempty"`
}

func (cc *CredentialController) FindByID(c *gin.Context) {
	var param struct{
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&param);err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, FindByIDResponse{
			Message: "validation error",
			Error: err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	credential, err := cc.credentialService.FindByID(ctx, param.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, FindByIDResponse{
			Message: "failed to find credential",
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, FindByIDResponse{
		Credential: credential,
		Message: "berhasil",
	})
}