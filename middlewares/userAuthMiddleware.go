package middlewares

import (
	"net/http"
	"strings"

	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/InnocentEdem/Go_Auth_v1/utils"
	"github.com/gin-gonic/gin"
)

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Malformed Authorization header"})
			c.Abort()
			return
		}

		claims, err := utils.ParseUserJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		var user models.ClientAppUser
		if err := initializers.DB.Where("email = ? AND client_id = ?", claims.Email, claims.ClientAppId).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found or does not belong to the client"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
