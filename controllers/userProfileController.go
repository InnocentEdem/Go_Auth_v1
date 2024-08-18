package controllers

import (
	"net/http"

	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Get user profile
// @Description Returns the user profile. Includes additional properties only if the client has use_additional_properties set to true. The id is not required in the request as it is derived from the authenticated user's context.
// @Tags user profile
// @Accept  json
// @Produce  json
// @Success 200 {object} UserProfileResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user/profile [get]
func GetUserProfile(c *gin.Context) {
	clientApp, exists := c.Get("clientApp")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Client information missing"})
		return
	}

	clientAppModel, ok := clientApp.(models.ClientApp)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving client information"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information missing"})
		return
	}

	userModel, ok := user.(models.ClientAppUser)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user information"})
		return
	}

	userProfile := UserProfileResponse{
		ID:        userModel.ID,
		FirstName: userModel.FirstName,
		LastName:  userModel.LastName,
		Email:     userModel.Email,
	}

	additionalProperties := AdditionalProperties{
		PhoneNumber:    userModel.AdditionalProperties.PhoneNumber,
		ProfilePicture: userModel.AdditionalProperties.ProfilePicture,
		DateOfBirth:    userModel.AdditionalProperties.DateOfBirth,
		Gender:         userModel.AdditionalProperties.Gender,
		LastLogin:      userModel.AdditionalProperties.LastLogin,
		Role:           userModel.AdditionalProperties.Role,
	}

	address := Address{
		Street:     userModel.AdditionalProperties.Address.Street,
		City:       userModel.AdditionalProperties.Address.City,
		State:      userModel.AdditionalProperties.Address.State,
		PostalCode: userModel.AdditionalProperties.Address.PostalCode,
		Country:    userModel.AdditionalProperties.Address.Country,
	}
	additionalProperties.Address = &address

	if clientAppModel.AppAdvancedConfig.UseAdditionalProperties {
		userProfile.AdditionalProperties = &additionalProperties
	}

	c.JSON(http.StatusOK, userProfile)
}

// @Summary Update user profile
// @Description Updates the user profile. Response includes additional properties only if the client has use_additional_properties set to true. The id is not required in the request as it is derived from the authenticated user's context.
// @Tags user profile
// @Accept  json
// @Produce  json
// @Param body body UpdateUserProfileRequest true "Update profile details"
// @Success 200 {object} UserProfileResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user/profile [put]
func UpdateUserProfile(c *gin.Context) {
	var body UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	clientApp, exists := c.Get("clientApp")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Client information missing"})
		return
	}

	clientAppModel, ok := clientApp.(models.ClientApp)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving client information"})
		return
	}
	userProps, exists := c.Get("clientAppUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information missing"})
		return
	}

	userModel, ok := userProps.(models.ClientAppUser)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user information"})
		return
	}

	var user models.ClientAppUser
	if err := initializers.DB.First(&user, "id = ?", userModel.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if body.FirstName != nil {

		user.FirstName = *body.FirstName
	}
	if body.LastName != nil {
		user.LastName = *body.LastName
	}
	if body.Email != nil {
		user.Email = *body.Email
	}

	if clientAppModel.AppAdvancedConfig.UseAdditionalProperties && body.AdditionalProperties != nil {
		user.AdditionalProperties = *body.AdditionalProperties
	}

	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user profile"})
		return
	}

	userProfile := UserProfileResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	additionalProperties := AdditionalProperties{
		PhoneNumber:    user.AdditionalProperties.PhoneNumber,
		ProfilePicture: user.AdditionalProperties.ProfilePicture,
		DateOfBirth:    user.AdditionalProperties.DateOfBirth,
		Gender:         user.AdditionalProperties.Gender,
		LastLogin:      user.AdditionalProperties.LastLogin,
		Role:           user.AdditionalProperties.Role,
	}

	address := Address{
		Street:     user.AdditionalProperties.Address.Street,
		City:       user.AdditionalProperties.Address.City,
		State:      user.AdditionalProperties.Address.State,
		PostalCode: user.AdditionalProperties.Address.PostalCode,
		Country:    user.AdditionalProperties.Address.Country,
	}
	additionalProperties.Address = &address

	if clientAppModel.AppAdvancedConfig.UseAdditionalProperties {
		userProfile.AdditionalProperties = &additionalProperties
	}

	c.JSON(http.StatusOK, userProfile)
}

type UserProfileResponse struct {
	ID                   uuid.UUID             `json:"id"`
	FirstName            string                `json:"first_name"`
	LastName             string                `json:"last_name"`
	Email                string                `json:"email"`
	AdditionalProperties *AdditionalProperties `json:"additional_properties,omitempty"`
}
type UpdateUserProfileRequest struct {
	FirstName            *string                      `json:"first_name,omitempty"`
	LastName             *string                      `json:"last_name,omitempty"`
	Email                *string                      `json:"email,omitempty"`
	AdditionalProperties *models.AdditionalProperties `json:"additional_properties,omitempty"`
}
