package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserSignup godoc
// @Summary User signup
// @Description Creates a new user account
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body UserSignupRequest true "Signup details"
// @Success 200 {object} UserSignupResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user/signup [post]
func UserSignup(c *gin.Context) {
	var body UserSignupRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid properties in request body"})
		return
	}

	clientApp, exists := c.Get("clientApp")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing client id"})
		return
	}

	clientAppModel, ok := clientApp.(models.ClientApp)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving client information"})
		return
	}

	var existingUser models.ClientAppUser
	if err := initializers.DB.Where("client_id = ? AND email = ?", clientAppModel.ID, body.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating account"})
		return
	}
	now := time.Now().UTC()

	user := models.ClientAppUser{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  string(hash),
		ClientAppID:  clientAppModel.ID,
	}

	if body.AdditionalProperties != nil {
		user.AdditionalProperties = models.AdditionalProperties{
			PhoneNumber:    body.AdditionalProperties.PhoneNumber,
			ProfilePicture: body.AdditionalProperties.ProfilePicture,
			DateOfBirth:    body.AdditionalProperties.DateOfBirth,
			Gender:         body.AdditionalProperties.Gender,
			Address: models.Address{
				Street:     body.AdditionalProperties.Address.Street,
				City:       body.AdditionalProperties.Address.City,
				State:      body.AdditionalProperties.Address.State,
				PostalCode: body.AdditionalProperties.Address.PostalCode,
				Country:    body.AdditionalProperties.Address.Country,
			},
			LastLogin: &now,
			Role:      body.AdditionalProperties.Role,
		}
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

type UserSignupRequest struct {
	FirstName            string                `json:"first_name" binding:"required"`
	LastName             string                `json:"last_name" binding:"required"`
	Email                string                `json:"email" binding:"required,email"`
	Password             string                `json:"password" binding:"required"`
	AdditionalProperties *AdditionalProperties `json:"additional_properties,omitempty"`
}

type AdditionalProperties struct {
	PhoneNumber    *string    `json:"phone_number,omitempty"`
	ProfilePicture *string    `json:"profile_picture,omitempty"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty"`
	Gender         *string    `json:"gender,omitempty"`
	Address        *Address   `json:"address,omitempty"`
	Role           *string    `json:"role,omitempty"`
	LastLogin      *time.Time `json:"last_login,omitempty"`
}

type Address struct {
	Street     *string `json:"street,omitempty"`
	City       *string `json:"city,omitempty"`
	State      *string `json:"state,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	Country    *string `json:"country,omitempty"`
}

func (a AdditionalProperties) MarshalJSON() ([]byte, error) {
	type Alias AdditionalProperties
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(a),
	})
}
