package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kominkamen/rootds-admin/models"
)

func GenerateJWT(user *models.AdminUser) (string, time.Time, error) {
	expirationTime := time.Now().AddDate(0, 2, 0)

	claims := models.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(accessTokenSecret)
	if err != nil {
		fmt.Println(err)
		return "", expirationTime, err
	}

	return tokenString, expirationTime, nil
}

func Sha256Hash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
