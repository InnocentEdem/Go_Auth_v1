package controllers

import (
	"net/http"

	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/InnocentEdem/Go_Auth_v1/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Refresh JWT
// @Description Refreshes the JWT and optionally the refresh token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   refresh_token  body  string  true  "Refresh Token"
// @Success 200 {object} RefreshTokenSuccessResponse
// @Failure 400 {object} RefreshTokenErrorResponse
// @Failure 401 {object} RefreshTokenErrorResponse
// @Failure 500 {object} RefreshTokenErrorResponse
// @Router /user/refresh-token [post]
func RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid properties in request body"})
		return
	}

	user, err := utils.ParseUserRefreshJWT(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	clientApp, exists := c.Get("clientApp")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing client id"})
		return
	}

	clientAppModel, ok := clientApp.(models.ClientApp)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving client information"})
		return
	}

	if !clientAppModel.AppAdvancedConfig.RefreshTokenEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token not activated for client"})
		return
	}

	var userData models.ClientAppUser
	if err := initializers.DB.Where("client_id = ? AND email = ?", clientAppModel.ID, user.Email).First(&userData).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		return
	}

	newToken, err := utils.GenerateUserJWT(userData, clientAppModel, "User")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new token"})
		return
	}
	newRefreshToken, err := utils.GenerateRefreshJWT(userData, clientAppModel, "User")
	if err != nil {
		c.JSON(http.StatusInternalServerError, RefreshTokenErrorResponse{Error: "Could not generate new refresh token"})
		return
	}

	response := RefreshTokenSuccessResponse{
		Message:      "Successful",
		LoginToken:   newToken,
		RefreshToken: newRefreshToken,
	}

	c.JSON(http.StatusOK, response)
}

type RefreshTokenSuccessResponse struct {
	Message      string `json:"message"`
	LoginToken   string `json:"login_token"`
	RefreshToken string `json:"refresh_token"`
}
type RefreshTokenErrorResponse struct {
	Error string `json:"error"`
}
