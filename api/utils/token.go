package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var privateKey = []byte(os.Getenv("JWT_KEY"))

func GenerateToken(user *models.User) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = user.ID
	claims["exp"] = now.AddDate(0, 0, 10).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(privateKey)
}

func ValidateToken(c *gin.Context) error {
	token, err := GetToken(c)

	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}

	return errors.New("Invalid token provided")
}

func CurrentUser(c *gin.Context) (models.User, error) {
	err := ValidateToken(c)
	if err != nil {
		return models.User{}, err
	}
	token, _ := GetToken(c)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	var user models.User
	dbs := database.DB.Where("id = ?", userId).Preload("Role").First(&user)
	if dbs.Error != nil {
		return models.User{}, dbs.Error
	}

	return user, nil
}

func GetToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

func getTokenFromRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
