package controllers

import (
	"net/http"

	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/InnocentEdem/Go_Auth_v1/services"
	"github.com/gin-gonic/gin"
)

type CreateClientAppRequest struct {
	AppName  string    `json:"app_name" binding:"required"`
}

func CreateNewClientApp(c *gin.Context) {
	var req CreateClientAppRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	client, exists := c.Get("client")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	clientModel, ok := client.(models.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving client information"})
		return
	}

	newClientApp := models.ClientApp{
		AppName:  req.AppName,
		ClientID: clientModel.ID,
	}

	err := services.CreateClientApp(newClientApp, clientModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ClientApp"})
		return
	}
	clientPtr, err := services.GetClientWithAppsByEmail(clientModel.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve client information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "App created successfully", "client": clientPtr})

}
