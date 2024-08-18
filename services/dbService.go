package services

import (
	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/InnocentEdem/Go_Auth_v1/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateClientAndRelatedEntities(client models.Client) error {
	return initializers.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&client).Error; err != nil {
			return err
		}

		defaultApp, err := CreateDefaultClientApp(tx, client.ID, "")
		if err != nil {
			return err
		}

		// Create default app configuration
		defaultConfig := utils.SetDefaultClientAppAdvancedConfig(defaultApp.ID)
		if err := tx.Create(&defaultConfig).Error; err != nil {
			return err
		}

		// Create default app confirmation method
		defaultConfirmationMethod := utils.SetDefaultClientAppConfirmationMethods(defaultApp.ID)
		if err := tx.Create(&defaultConfirmationMethod).Error; err != nil {
			return err
		}
		return nil
	})
}

func GetClientByEmail(email string) (*models.Client, error) {
	var client models.Client
	if err := initializers.DB.Where("email = ?", email).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func GetClientWithAppsByEmail(email string) (*models.Client, error) {
	var client models.Client
	if err := initializers.DB.Preload("ClientApps").Preload("ClientApps.AppAdvancedConfig").Preload("ClientApps.AppConfirmationMethods").
		Where("email = ?", email).
		First(&client).Error; err != nil {
		return nil, err
	}

	return &client, nil
}

func GetClientByID(id uuid.UUID) (*models.Client, error) {
	var client models.Client
	if err := initializers.DB.First(&client, id).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func GetClientAppByID(id uuid.UUID) (*models.ClientApp, error) {
	var app models.ClientApp
	if err := initializers.DB.First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func GetAppAdvancedConfigByID(id uuid.UUID) (*models.AppAdvancedConfig, error) {
	var config models.AppAdvancedConfig
	if err := initializers.DB.First(&config, id).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// GetClientAppsByClientID fetches all client apps by client ID
func GetClientAppsByClientID(clientID uuid.UUID) ([]models.ClientApp, error) {
	var apps []models.ClientApp
	if err := initializers.DB.Where("client_id = ?", clientID).Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func GetAppConfirmationMethodByID(id uuid.UUID) (*models.AppConfirmationMethod, error) {
	var method models.AppConfirmationMethod
	if err := initializers.DB.First(&method, id).Error; err != nil {
		return nil, err
	}
	return &method, nil
}

func CreateClient(firstName, lastName, email, password string) (*models.Client, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	client := models.Client{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hash),
	}

	if err := initializers.DB.Create(&client).Error; err != nil {
		return nil, err
	}

	return &client, nil
}

func CreateDefaultClientApp(tx *gorm.DB, clientID uuid.UUID, name string) (*models.ClientApp, error) {
	if name == "" {
		name = "default_app"
	}

	apn, err := utils.GenerateAPN(16)
	if err != nil {
		return nil, err
	}

	defaultApp := models.ClientApp{
		AppName:  name,
		ClientID: clientID,
		APN:      apn,
	}

	if err := tx.Create(&defaultApp).Error; err != nil {
		return nil, err
	}

	return &defaultApp, nil
}

func CreateDefaultAppConfig(clientAppID uuid.UUID) (*models.AppAdvancedConfig, error) {
	defaultAppConfig := utils.SetDefaultClientAppAdvancedConfig(clientAppID)

	if err := initializers.DB.Create(&defaultAppConfig).Error; err != nil {
		return nil, err
	}

	return &defaultAppConfig, nil
}

func CreateDefaultAppConfirmationMethod(appID uuid.UUID) (*models.AppConfirmationMethod, error) {
	defaultAppConfirmationMethod := utils.SetDefaultClientAppConfirmationMethods(appID)

	if err := initializers.DB.Create(&defaultAppConfirmationMethod).Error; err != nil {
		return nil, err
	}

	return &defaultAppConfirmationMethod, nil
}

func CreateClientApp(clientApp models.ClientApp, clientId uuid.UUID) error {

	return initializers.DB.Transaction(func(tx *gorm.DB) error {

		newApp, err := CreateDefaultClientApp(tx, clientId, clientApp.AppName)
		if err != nil {
			return err
		}

		defaultConfig := utils.SetDefaultClientAppAdvancedConfig(newApp.ID)
		if err := tx.Create(&defaultConfig).Error; err != nil {
			return err
		}

		defaultConfirmationMethod := utils.SetDefaultClientAppConfirmationMethods(newApp.ID)
		if err := tx.Create(&defaultConfirmationMethod).Error; err != nil {
			return err
		}
		return nil
	})
}

type AppWithUsers struct {
	AppName string
	AppID   uuid.UUID
	Users   []models.ClientAppUser
}

func GetUsersByAppID(appID uuid.UUID) ([]models.ClientAppUser, error) {
	var users []models.ClientAppUser

	err := initializers.DB.Where("client_app_id = ?", appID).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetAllClientApps(clientID uuid.UUID) ([]models.ClientApp, error) {
	var clientApps []models.ClientApp

	err := initializers.DB.
		Preload("AppAdvancedConfig").
		Preload("AppConfirmationMethods").
		Where("client_id = ?", clientID).
		Find(&clientApps).Error
	if err != nil {
		return nil, err
	}

	return clientApps, nil
}

func UpdateAppAdvancedConfig(clientAppID uuid.UUID, clientID uuid.UUID, updates models.AppAdvancedConfig) (*models.AppAdvancedConfig, error) {
    var appConfig models.AppAdvancedConfig

    err := initializers.DB.
        Joins("JOIN client_apps ON client_apps.id = app_advanced_configs.client_app_id").
        Where("client_apps.client_id = ? AND client_apps.id = ?", clientID, clientAppID).
        First(&appConfig).Error

    if err != nil {
        if gorm.ErrRecordNotFound == err {
            return nil, nil 
        }
        return nil, err
    }

    updatesMap := map[string]interface{}{}

    if updates.CorsAllowedOrigins != nil {
        updatesMap["cors_allowed_origins"] = updates.CorsAllowedOrigins
    }
    if updates.JWTExpiryTime != 0 {
        updatesMap["jwt_expiry_time"] = updates.JWTExpiryTime
    }
    if updates.RefreshTokenEnabled {
        updatesMap["refresh_token_enabled"] = updates.RefreshTokenEnabled
    }
    if updates.RefreshTokenExpiryTime != 0 {
        updatesMap["refresh_token_expiry_time"] = updates.RefreshTokenExpiryTime
    }
    if updates.AllowJWTCustomClaims {
        updatesMap["allow_jwt_custom_claims"] = updates.AllowJWTCustomClaims
    }
    if updates.UseAdditionalProperties {
        updatesMap["use_additional_properties"] = updates.UseAdditionalProperties
    }

    err = initializers.DB.Model(&appConfig).Updates(updatesMap).Error
    if err != nil {
        return nil, err
    }

    return &appConfig, nil
}

