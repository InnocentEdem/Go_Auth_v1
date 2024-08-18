package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientApp struct {
	ID                     uuid.UUID             `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AppName                string                `gorm:"type:varchar(100);not null" json:"app_name"`
	ClientID               uuid.UUID             `gorm:"type:uuid;not null" json:"client_id"`
	CreatedAt              time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt              gorm.DeletedAt        `gorm:"index" json:"deleted_at,omitempty"`
	APN                    string                `gorm:"type:varchar(100);uniqueIndex" json:"apn"`
	AppAdvancedConfig      AppAdvancedConfig     `gorm:"foreignKey:ClientAppID;constraint:OnDelete:CASCADE" json:"app_advanced_config"`
	AppConfirmationMethods AppConfirmationMethod `gorm:"foreignKey:ClientAppID;constraint:OnDelete:CASCADE" json:"app_confirmation_methods"`
}

func (app *ClientApp) BeforeCreate(tx *gorm.DB) (err error) {
	app.ID = uuid.New()

	if app.APN, err = GenerateAPN(16); err != nil {
		return err
	}

	return nil
}

func GenerateAPN(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
