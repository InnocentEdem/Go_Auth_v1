package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"authapp.com/m/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type ClientClaims struct {
    Email string `json:"email"`
	Role string `json:"role"`
    FirstName string `json:"firstname"`
    LastName string `json:"lastname"`
    Apn string `json:"apn"`
    jwt.RegisteredClaims
}

type UserClaims struct {
	Email string `json:"email"`
	Role string `json:"role"`
    FirstName string `json:"firstname"`
    LastName string `json:"lastname"`
	ClientId uuid.UUID `json:"client"`
    jwt.RegisteredClaims
}
var jwtKey = []byte(os.Getenv("JWT_KEY"))


func GenerateJWT(client models.Client, role string) (string, error) {
    expirationTime := time.Now().Add(2 * time.Hour) 
    claims := &ClientClaims{
        Email: client.Email,
		Role: role,
        FirstName: client.FirstName,
        LastName: client.LastName,
        Apn: client.APN,
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

func GenerateUserJWT(user models.User, role string) (string, error) {
    expirationTime := time.Now().Add(2 * time.Hour) 
    claims := &UserClaims{
        Email: user.Email,
		Role: role,
        FirstName: user.FirstName,
        LastName: user.LastName,
		ClientId: user.ClientID,
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