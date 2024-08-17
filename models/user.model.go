package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                   uuid.UUID            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FirstName            string               `gorm:"type:varchar(100);not null"`
	LastName             string               `gorm:"type:varchar(100);not null"`
	Email                string               `gorm:"type:varchar(100);not null"`
	Password             string               `gorm:"type:varchar(255);not null"`
	AdditionalProperties AdditionalProperties `gorm:"embedded"`
	UserConfirmation     UserConfirmation     `gorm:"constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ClientID  uuid.UUID      `gorm:"type:uuid"`
	Client    Client
}

type AdditionalProperties struct {
	PhoneNumber    *string    `gorm:"type:varchar(20)"`
	ProfilePicture *string    `gorm:"type:varchar(255)"`
	DateOfBirth    *time.Time `gorm:"type:date"`
	Gender         *string    `gorm:"type:varchar(20)"`
	Address        Address    `gorm:"embedded"`
	LastLogin      *time.Time `gorm:"type:timestamp"`
	Role           *string    `gorm:"type:varchar(20)"`
}

type Address struct {
	Street     *string `gorm:"type:varchar(255)"`
	City       *string `gorm:"type:varchar(100)"`
	State      *string `gorm:"type:varchar(100)"`
	PostalCode *string `gorm:"type:varchar(20)"`
	Country    *string `gorm:"type:varchar(100)"`
}

func (User) TableName() string {
	return "users"
}
