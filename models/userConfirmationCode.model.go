package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type ConfirmationCode struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    UserID    uuid.UUID      `gorm:"type:uuid;unique" json:"user_id"`
    Code      string         `gorm:"size:6;not null" json:"code"`
    ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

