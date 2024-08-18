package middlewares

import (
	"net/http"
	"strings"

	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/gin-gonic/gin"
)

func DynamicCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		client, exists := c.Get("client")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing client id"})
			return
		}

		clientModel, ok := client.(models.Client)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving client information"})
			return
		}

		var config models.AppAdvancedConfig
		if err := initializers.DB.Where("client_id = ?", clientModel.ID).First(&config).Error; err != nil {
			c.Next()
			return
		}

		if len(config.CorsAllowedOrigins) > 0 && config.CorsAllowedOrigins[0] != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(config.CorsAllowedOrigins, ","))
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
