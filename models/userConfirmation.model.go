package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserConfirmation struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID                 uuid.UUID      `gorm:"type:uuid;unique" json:"user_id"`
	ClientID               uuid.UUID      `gorm:"type:uuid" json:"client_id"`
	EmailConfirmed         bool           `gorm:"default:false" json:"email_confirmed"`
	PhoneConfirmed         bool           `gorm:"default:false" json:"phone_confirmed"`
	PaymentMethodConfirmed bool           `gorm:"default:false" json:"payment_method_confirmed"`
	CreatedAt              time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

