package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type TClaims struct {
	UserID  int32  `json:"user_id"`
	Account string `json:"account"`
	jwt.RegisteredClaims
}

func GetToken(account string, userid int32) (string, error) {
	claims := TClaims{
		userid,
		account,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(loadSecretKey())
}

func DecodeToken(tokenValue string) (int32, string, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &TClaims{}, func(token *jwt.Token) (interface{}, error) {
		return loadSecretKey(), nil
	})
	if err != nil {
		return 0, "", err
	}
	claims, ok := token.Claims.(*TClaims)
	if ok && token.Valid {
		return claims.UserID, claims.Account, nil
	}
	return 0, "", fmt.Errorf("invalid token %s", tokenValue)
}

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		//c.AbortWithStatus(http.StatusUnauthorized)
		//c.JSON(http.StatusOK, Failure("Unauthorized"))

		return
	}

	tokenString := authHeader[len("Bearer "):]
	userID, account, err := DecodeToken(tokenString)
	if err != nil {
		//c.AbortWithStatus(http.StatusUnauthorized)
		c.JSON(http.StatusUnauthorized, Failure("Unauthorized"))
		return
	}

	c.Set("userid", userID)
	c.Set("account", account)
	c.Next()
}

func loadSecretKey() []byte {
	secretKey := []byte(GetDefaultKey())
	return secretKey
}
