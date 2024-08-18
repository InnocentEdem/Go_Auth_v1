package utils

import (
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/google/uuid"
)

func SetDefaultApp(clientID uuid.UUID) models.ClientApp {
	return models.ClientApp{
		AppName:  "Default App",
		ClientID: clientID,
	}
}
