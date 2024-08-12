package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
    gorm.Model
    ID                     uuid.UUID           `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    FirstName              string              `gorm:"type:varchar(100);not null"`
    LastName               string              `gorm:"type:varchar(100);not null"`
    Email                  string              `gorm:"type:varchar(100);uniqueIndex;not null"`
    Password               string              `gorm:"type:varchar(255);not null"`
    APN                    string              `gorm:"type:varchar(100);uniqueIndex"`
    ClientAdvancedConfig   ClientAdvancedConfig `gorm:"foreignKey:ClientID"`
    CreatedAt              time.Time
    UpdatedAt              time.Time
    DeletedAt              gorm.DeletedAt      `gorm:"index"`
    Users                  []User
    Apps                   []App               `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE;"`
}


func (c *Client) BeforeCreate(tx *gorm.DB) (err error) {
    c.ID = uuid.New()

    if c.APN, err = GenerateAPN(16); err != nil {
        return err
    }

    return nil
}


func GenerateAPN(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", fmt.Errorf("failed to generate random bytes: %w", err)
    }
    return hex.EncodeToString(bytes), nil
}
