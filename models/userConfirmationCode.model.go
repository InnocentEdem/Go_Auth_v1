package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type ConfirmationCode struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    UserID    uuid.UUID      `gorm:"type:uuid;unique"`
    Code      string         `gorm:"size:6;not null"`
    ExpiresAt time.Time      `gorm:"not null"`  
    CreatedAt time.Time      `gorm:"autoCreateTime"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
