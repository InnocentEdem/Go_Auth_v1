package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientAppUser struct {
	ID                   uuid.UUID            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FirstName            string               `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName             string               `gorm:"type:varchar(100);not null" json:"last_name"`
	Email                string               `gorm:"type:varchar(100);not null" json:"email"`
	Password             string               `gorm:"type:varchar(255);not null" json:"password"`
	AdditionalProperties AdditionalProperties `gorm:"embedded" json:"additional_properties"`
	UserConfirmation     UserConfirmation     `gorm:"foreignKey:UserConfirmationID;constraint:OnDelete:CASCADE" json:"user_confirmation"`
	UserConfirmationID   uuid.UUID            `gorm:"type:uuid;not null" json:"user_confirmation_id"`
	CreatedAt            time.Time            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time            `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt            gorm.DeletedAt       `gorm:"index" json:"deleted_at,omitempty"`
	ClientAppID          uuid.UUID            `gorm:"type:uuid;index;not null" json:"client_app_id"`
}

type AdditionalProperties struct {
	PhoneNumber    *string    `gorm:"type:varchar(20)" json:"phone_number,omitempty"`
	ProfilePicture *string    `gorm:"type:varchar(255)" json:"profile_picture,omitempty"`
	DateOfBirth    *time.Time `gorm:"type:date" json:"date_of_birth,omitempty"`
	Gender         *string    `gorm:"type:varchar(20)" json:"gender,omitempty"`
	Address        Address    `gorm:"embedded" json:"address"`
	LastLogin      *time.Time `gorm:"type:timestamp" json:"last_login,omitempty"`
	Role           *string    `gorm:"type:varchar(20)" json:"role,omitempty"`
}

type Address struct {
	Street     *string `gorm:"type:varchar(255)" json:"street,omitempty"`
	City       *string `gorm:"type:varchar(100)" json:"city,omitempty"`
	State      *string `gorm:"type:varchar(100)" json:"state,omitempty"`
	PostalCode *string `gorm:"type:varchar(20)" json:"postal_code,omitempty"`
	Country    *string `gorm:"type:varchar(100)" json:"country,omitempty"`
}
