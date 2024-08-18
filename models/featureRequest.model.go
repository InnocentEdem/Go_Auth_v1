package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeatureRequest struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Feature   string         `gorm:"type:varchar(100);not null" json:"feature"`
	Title     string         `gorm:"type:varchar(100);not null" json:"title"`
	Email     string         `gorm:"type:varchar(100);not null" json:"email"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (u *ClientAppUser) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
