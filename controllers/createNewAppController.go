package controllers

import (
	"net/http"

	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateNewApp(c *gin.Context) {
	var body CreateNewAppRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	client, exists := c.Get("client")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing client id"})
		return
	}

	clientModel, ok := client.(models.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error retrieving client information"})
		return
	}

	newApp := models.ClientApp{
		AppName:  body.AppName,
		ClientID: clientModel.ID,
	}

	if err := initializers.DB.Create(&newApp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create app"})
		return
	}
	response := CreateNewAppResponse{
		Message: "App created successfully",
		App: AppDetails{
			ID:       newApp.ID,
			AppName:  newApp.AppName,
			ClientID: newApp.ClientID,
		},
	}

	c.JSON(http.StatusOK, response)
}

type CreateNewAppRequest struct {
	AppName string `json:"app_name" binding:"required"`
}
type CreateNewAppResponse struct {
	Message string     `json:"message"`
	App     AppDetails `json:"app"`
}

type AppDetails struct {
	ID        uuid.UUID `json:"id"`
	AppName   string    `json:"app_name"`
	ClientID  uuid.UUID `json:"client_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
