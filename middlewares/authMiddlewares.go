package middlewares

import (
	"lamhat/core"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type MyClaim struct {
	UserId int `json:"Subject"`
	jwt.RegisteredClaims
}

func ValidateAuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("authToken")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, "No access - Token is missing")
			c.Abort()
			return
		}

		// Parse the token
		token, err := jwt.ParseWithClaims(tokenString, &MyClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(core.Config.JWT.SECRET_KEY), nil
		})

		if err != nil {
			core.Sugar.Errorf("Error parsing token: %s", err.Error())
			c.JSON(http.StatusUnauthorized, "No access - Invalid token")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*MyClaim); ok && token.Valid {
			// Token is valid, you can use `claims.UserId` if needed
			c.Set("userId", claims.UserId) // Example of storing user ID in context
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, "No access - Invalid token")
			c.Abort()
		}
	}
}
