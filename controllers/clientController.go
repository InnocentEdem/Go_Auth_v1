package controllers

import (
	"fmt"
	"net/http"

	"github.com/InnocentEdem/Go_Auth_v1/initializers"
	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/InnocentEdem/Go_Auth_v1/services"
	"github.com/InnocentEdem/Go_Auth_v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func ClientSignup(c *gin.Context) {
	var body struct {
		FirstName string `json:"firstname" binding:"required"`
		LastName  string `json:"lastname" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid properties in request body"})
		return
	}

	client, err := services.GetClientByEmail(body.Email)
	if err == nil && client != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with email already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating account"})
		return
	}

	var newClient models.Client = models.Client{
        FirstName: body.FirstName,
        LastName:  body.LastName,
        Email:     body.Email,
        Password:  string(hash),
    }
    
    err = services.CreateClientAndRelatedEntities(newClient)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign up client"})
        return
    }
    
  
    newClientPtr, err := services.GetClientWithAppsByEmail(body.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve client information"})
        return
    }




	c.JSON(http.StatusOK, gin.H{"message": "Client created successfully", "client": newClientPtr})
}


func ClientLogin(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid properties in request body"})
		return
	}

	var client models.Client
	if err := initializers.DB.Where("email = ?", body.Email).First(&client).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}
	token, err := utils.GenerateJWT(client, "CLIENT")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "login_token": token})
}

func ClientUpdatePassword(c *gin.Context) {
	var request ClientPasswordUpdateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Missing or invalid properties in request body"})
		return
	}

	client, exists := c.Get("client")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Client not authenticated"})
		return
	}

	clientModel, ok := client.(models.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error retrieving client information"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(clientModel.Password), []byte(request.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Incorrect old password"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error updating password"})
		return
	}

	clientModel.Password = string(hash)
	if err := initializers.DB.Save(&clientModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error updating password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update successful"})
}

func GetUsersByClient(c *gin.Context) {
	client, exists := c.Get("client")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	clientModel, ok := client.(models.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error retrieving client information"})
		return
	}

	var users []models.ClientAppUser
	if err := initializers.DB.Where("client_id = ?", clientModel.ID).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error fetching users"})
		return
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		})
	}

	c.JSON(http.StatusOK, userResponses)
}

func DeleteUserByClient(c *gin.Context) {
	var request DeleteUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	client, exists := c.Get("client")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	clientModel, ok := client.(models.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error retrieving client information"})
		return
	}

	var user models.ClientAppUser
	if err := initializers.DB.Where("id = ? AND client_id = ?", request.UserID, clientModel.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User not found or does not belong to this client"})
		return
	}

	if err := initializers.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, DeleteUserResponse{Message: "User deleted successfully"})
}

func GetClientAppAPN(c *gin.Context) {

	var requestBody struct {
		AppID uuid.UUID `json:"app_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}
	client, exists := c.Get("client")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	clientModel, ok := client.(models.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error retrieving client information"})
		return
	}

	var clientApp models.ClientApp
	if err := initializers.DB.Where("id = ? AND client_id = ?", requestBody.AppID, clientModel.ID).First(&clientApp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error fetching client app from database"})
		return
	}

	c.JSON(http.StatusOK, GetClientAPNResponse{APN: clientApp.APN})
}

func HandleFeatureRequest(c *gin.Context) {
	var request FeatureRequestInput

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	client, exists := c.Get("client")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	clientModel, ok := client.(models.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error retrieving client information"})
		return
	}

	featureRequest := utils.FeatureRequest{
		FeatureName:        request.Title,
		FeatureDescription: request.Feature,
		SenderName:         fmt.Sprintf("%s %s", clientModel.FirstName, clientModel.LastName),
		SenderEmail:        clientModel.Email,
	}

	if err := utils.SendFeatureRequestEmail(featureRequest); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feature request email sent successfully"})
}

func GetClient(c *gin.Context) {
	client, exists := c.Get("client")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Client information missing"})
		return
	}

	clientModel, ok := client.(models.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving client information"})
		return
	}
    var response models.Client
    if err := initializers.DB.Preload("ClientApp").Where("email = ?", clientModel.Email).First(&client).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid client"})
        c.Abort()
        return
    }

	c.JSON(http.StatusOK, response)
}

func GetUsersByAppID(c *gin.Context) {
    var req AppIDRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    users, err := services.GetUsersByAppID(req.AppID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"app_id": req.AppID, "users": users})
}

func GetAllClientApps(c *gin.Context) {

    client, exists := c.Get("client")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Client information missing"})
		return
	}

	clientModel, ok := client.(models.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving client information"})
		return
	}
    clientApps, err := services.GetAllClientApps(clientModel.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch client apps"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"client_apps": clientApps})
}

type AppIDRequest struct {
    AppID uuid.UUID `json:"app_id" binding:"required"`
}

type ClientPasswordUpdateRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}
type DeleteUserRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type GetClientAPNResponse struct {
	APN string `json:"apn"`
}

type FeatureRequest struct {
	FeatureName        string `json:"feature_name" binding:"required"`
	FeatureDescription string `json:"feature_description" binding:"required"`
}

type ClientResponse struct {
	ID                   uuid.UUID                    `json:"id"`
	FirstName            string                       `json:"firstname"`
	LastName             string                       `json:"lastname"`
	ClientAdvancedConfig ClientAdvancedConfigResponse `json:"client_advanced_config"`
}
