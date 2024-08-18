package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/InnocentEdem/Go_Auth_v1/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type ClientClaims struct {
	Email     string `json:"email"`
	Role      string `json:"role"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	jwt.RegisteredClaims
}

type UserClaims struct {
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	ClientAppId  uuid.UUID `json:"client_app_id"`
	jwt.RegisteredClaims
}
type UserRefreshTokenClaims struct {
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	ClientAppId  uuid.UUID `json:"client_app_id"`
	Type      string    `json:"type"`
	jwt.RegisteredClaims
}

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func GenerateJWT(client models.Client, role string) (string, error) {
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &ClientClaims{
		Email:     client.Email,
		Role:      role,
		FirstName: client.FirstName,
		LastName:  client.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	fmt.Println(os.Getenv("JWT_KEY"), err)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateUserJWT(user models.ClientAppUser, clientApp models.ClientApp, role string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(clientApp.AppAdvancedConfig.JWTExpiryTime) * time.Second)
	claims := &UserClaims{
		Email:     user.Email,
		Role:      role,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		ClientAppId:  user.ClientAppID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWTWithClaims(tokenString string) (*ClientClaims, error) {
	claims := &ClientClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ClientClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

func ParseUserJWT(tokenString string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

func GenerateRefreshJWT(user models.ClientAppUser, clientApp models.ClientApp, role string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(int64(clientApp.AppAdvancedConfig.RefreshTokenExpiryTime)) * time.Second)
	claims := &UserRefreshTokenClaims{
		Email:     user.Email,
		Role:      role,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		ClientAppId:  user.ClientAppID,
		Type:      "refresh_token",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseUserRefreshJWT(tokenString string) (*UserRefreshTokenClaims, error) {
	claims := &UserRefreshTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtKey, nil
	})
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserRefreshTokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
