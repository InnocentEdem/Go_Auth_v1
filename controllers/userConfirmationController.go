package controllers

import (
	"net/http"
	"time"

	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/gin-gonic/gin"
)

func ValidateConfirmationCode(c *gin.Context) {
	var request struct {
		Code string `json:"validation_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userModel, ok := user.(*models.ClientAppUser)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var confirmationCode models.ConfirmationCode

	if err := initializers.DB.Where("user_id = ? AND code = ?", userModel.ID, request.Code).First(&confirmationCode).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User or code does not exist"})
		return
	}

	if time.Now().After(confirmationCode.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code has expired"})
		return
	}

	if err := initializers.DB.Delete(&confirmationCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Confirmation successful"})
}
