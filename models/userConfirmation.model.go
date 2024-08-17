package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserConfirmation struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID                 uuid.UUID      `gorm:"type:uuid;unique"`
	ClientID               uuid.UUID      `gorm:"type:uuid"`
	EmailConfirmed         bool           `gorm:"default:false"`
	PhoneConfirmed         bool           `gorm:"default:false"`
	PaymentMethodConfirmed bool           `gorm:"default:false"`
	CreatedAt              time.Time      `gorm:"autoCreateTime"`
	UpdatedAt              time.Time      `gorm:"autoUpdateTime"`
	DeletedAt              gorm.DeletedAt `gorm:"index"`
}
