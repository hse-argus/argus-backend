package middleware

import (
	"argus-backend/internal/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

func EnableCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "https://argus.appweb.space")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.Writer.WriteHeader(http.StatusOK)
		return
	}

	c.Next()
}

type Claims struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Id    int    `json:"id"`
	jwt.RegisteredClaims
}

func JWTTokenVerify(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		logger.Error("token is empty")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
		c.Abort()
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		jwtKey, _ := os.LookupEnv("SECRET_KEY")
		return []byte(jwtKey), nil
	})

	if err != nil || !token.Valid {
		logger.Error("token is invalid")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid"})
		c.Abort()
		return
	}

	c.Set("login", claims.Login)
	c.Set("name", claims.Name)
	c.Set("id", claims.Id)

	c.Next()
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		jwtKey, _ := os.LookupEnv("SECRET_KEY")
		return []byte(jwtKey), nil
	})

	if err != nil || !token.Valid {
		logger.Error("token is invalid")
		return nil, err
	}

	return claims, nil
}
