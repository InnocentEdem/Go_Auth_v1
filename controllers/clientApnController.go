package controllers

import (
	"net/http"

	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/InnocentEdem/Go_Auth_v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GenerateAPN(c *gin.Context) {
    var requestBody struct {
		AppID uuid.UUID `json:"app_id" binding:"required"`
	}
    if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
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

    var clientAppModel models.ClientApp
	if err := initializers.DB.Where("client_id = ? AND id = ?", clientModel.ID, requestBody.AppID).First(&clientAppModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "ClientApp not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching ClientApp"})
		}
		return
	}

    apn, err := utils.GenerateAPN(16)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate APN"})
        return
    }


    clientAppModel.APN = string(apn)

    if err := initializers.DB.Save(&clientAppModel).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update APN"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"apn": apn})
}