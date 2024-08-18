package controllers

import (
	"net/http"

	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/InnocentEdem/Go_Auth_v1/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UpdateClientAdvancedConfigHandler(c *gin.Context) {
	clientAppID := c.Param("id")
	if clientAppID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AppId ID is required"})
		return
	}
	parsedAppConfigID, err := uuid.Parse(clientAppID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid AppAdvancedConfig ID"})
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

	var updates models.AppAdvancedConfig
	if err := c.BindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedAppConfig, err := services.UpdateAppAdvancedConfig(parsedAppConfigID,clientModel.ID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update AppAdvancedConfig"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated_config": updatedAppConfig})
}

func GetClientAdvancedConfig(c *gin.Context) {

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

	var clientConfig models.AppAdvancedConfig
	if err := initializers.DB.First(&clientConfig, "client_id = ?", clientModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client configuration not found"})
		return
	}

	response := ClientAdvancedConfigResponse{
		ID:                      clientConfig.ID,
		ClientAppID:             clientConfig.ClientAppID,
		CorsAllowedOrigins:      clientConfig.CorsAllowedOrigins,
		JWTExpiryTime:           clientConfig.JWTExpiryTime,
		RefreshTokenEnabled:     clientConfig.RefreshTokenEnabled,
		RefreshTokenExpiryTime:  clientConfig.RefreshTokenExpiryTime,
		AllowJWTCustomClaims:    clientConfig.AllowJWTCustomClaims,
		UseAdditionalProperties: clientConfig.UseAdditionalProperties,
	}

	c.JSON(http.StatusOK, response)
}

type ClientAdvancedConfigResponse struct {
	ID                      uuid.UUID `json:"id"`
	ClientAppID             uuid.UUID `json:"client_app_id"`
	CorsAllowedOrigins      []string  `json:"cors_allowed_origins"`
	JWTExpiryTime           int       `json:"jwt_expiry_time"`
	RefreshTokenEnabled     bool      `json:"refresh_token_enabled"`
	RefreshTokenExpiryTime  int       `json:"refresh_token_expiry_time"`
	AllowJWTCustomClaims    bool      `json:"allow_jwt_custom_claims"`
	UseAdditionalProperties bool      `json:"use_additional_properties"`
}
type UpdateClientAdvancedConfigRequest struct {
	CorsAllowedOrigins      *[]string `json:"cors_allowed_origins,omitempty"`
	JWTExpiryTime           *int      `json:"jwt_expiry_time,omitempty"`
	RefreshTokenEnabled     *bool     `json:"refresh_token_enabled,omitempty"`
	RefreshTokenExpiryTime  *int      `json:"refresh_token_expiry_time,omitempty"`
	AllowJWTCustomClaims    *bool     `json:"allow_jwt_custom_claims,omitempty"`
	UseAdditionalProperties *bool     `json:"use_additional_properties,omitempty"`
}
