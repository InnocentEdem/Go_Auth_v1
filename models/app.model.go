package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type App struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    AppName     string    `gorm:"type:varchar(100);not null"`
    AppPublicId string    `gorm:"type:varchar(100);not null;uniqueIndex"`
    ClientID    uuid.UUID `gorm:"type:uuid;not null"`
    Client      Client    `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE;"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}
