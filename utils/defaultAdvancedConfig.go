package utils

import (
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func SetDefaultClientAppAdvancedConfig(clientAppId uuid.UUID) models.AppAdvancedConfig {
	return models.AppAdvancedConfig{
		ClientAppID:             clientAppId,
		CorsAllowedOrigins:      pq.StringArray{""},
		JWTExpiryTime:           3600,
		RefreshTokenEnabled:     false,
		RefreshTokenExpiryTime:  7200,
		AllowJWTCustomClaims:    false,
		UseAdditionalProperties: false,
	}
}
func SetDefaultClientAppConfirmationMethods(clientAppId uuid.UUID) models.AppConfirmationMethod {
	return models.AppConfirmationMethod{
		ClientAppID:          clientAppId,
		ConfirmEmail:         false,
		ConfirmPhone:         false,
		ConfirmPaymentMethod: false,
	}
}
