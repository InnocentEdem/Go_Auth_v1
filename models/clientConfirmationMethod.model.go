package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppConfirmationMethod struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ClientAppID          uuid.UUID      `gorm:"type:uuid;unique" json:"client_app_id"`
	ConfirmEmail         bool           `gorm:"default:false" json:"confirm_email"`
	ConfirmPhone         bool           `gorm:"default:false" json:"confirm_phone"`
	ConfirmPaymentMethod bool           `gorm:"default:false" json:"confirm_payment_method"`
	CreatedAt            time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

