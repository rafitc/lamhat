package routes

import (
	"lamhat/core"
	"lamhat/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Sugar_logger = core.Sugar

func AddUserRoutes(rg *gin.RouterGroup) {

	user := rg.Group("/")
	user.Use(middlewares.ValidateAuthToken())

	user.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "User Handle")
	})
	user.GET("/home", func(c *gin.Context) {
		c.JSON(http.StatusOK, "User Home")
	})
}
