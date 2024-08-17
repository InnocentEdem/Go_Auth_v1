package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientConfirmationMethod struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ClientID             uuid.UUID      `gorm:"type:uuid;unique"`
	ConfirmEmail         bool           `gorm:"default:false"`
	ConfirmPhone         bool           `gorm:"default:false"`
	ConfirmPaymentMethod bool           `gorm:"default:false"`
	CreatedAt            time.Time      `gorm:"autoCreateTime"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}
