package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FirstName   string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName    string         `gorm:"type:varchar(100);not null" json:"last_name"`
	Email       string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password    string         `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	ClientApps  []ClientApp    `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE;" json:"client_apps"`
}

