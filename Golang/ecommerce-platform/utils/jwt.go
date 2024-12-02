// utils/jwt.go
package utils

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
    UserID int    `json:"user_id"`
    Role   string `json:"role"`
    jwt.StandardClaims
}

func GenerateJWT(userID int, role string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: userID,
        Role:   role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }
    log.Printf("Generating token for UserID: %d, Role: %s", userID, role)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        return nil, err
    }
    if !token.Valid {
        return nil, err
    }
    return claims, nil
}
