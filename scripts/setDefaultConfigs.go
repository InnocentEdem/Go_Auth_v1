package scripts

import (
	"log"

	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/InnocentEdem/Go_Auth_v1/utils"
	"gorm.io/gorm"
)

func SetDefaultConfigScript() {

	var clients []models.Client
	if err := initializers.DB.Find(&clients).Error; err != nil {
		log.Fatalf("failed to fetch clients: %v", err)
	}

	for _, client := range clients {
		var config models.AppAdvancedConfig
		if err := initializers.DB.Where("client_id = ?", client.ID).First(&config).Error; err == gorm.ErrRecordNotFound {
			defaultConfig := utils.SetDefaultClientAppAdvancedConfig(client.ID)
			if err := initializers.DB.Create(&defaultConfig).Error; err != nil {
				log.Printf("failed to create default config for client %s: %v", client.ID, err)
			} else {
				log.Printf("default config created for client %s", client.ID)
			}
		}
	}
}
