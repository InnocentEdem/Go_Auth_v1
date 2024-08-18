package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AppAdvancedConfig struct {
	ID                      uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ClientAppID             uuid.UUID      `gorm:"type:uuid;unique" json:"client_app_id"`
	CorsAllowedOrigins      pq.StringArray `gorm:"type:text[]" json:"cors_allowed_origins"`
	JWTExpiryTime           int            `json:"jwt_expiry_time"`
	RefreshTokenEnabled     bool           `json:"refresh_token_enabled"`
	RefreshTokenExpiryTime  int            `json:"refresh_token_expiry_time"`
	AllowJWTCustomClaims    bool           `json:"allow_jwt_custom_claims"`
	UseAdditionalProperties bool           `json:"use_additional_properties"`
	CreatedAt               time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}


func (Client) TableName() string {
	return "clients"
}
