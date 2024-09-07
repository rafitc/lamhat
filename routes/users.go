package routes

import (
	"lamhat/core"
	"lamhat/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Sugar_logger = core.Sugar

func AddUserRoutes(rg *gin.RouterGroup) {

	users := rg.Group("/")
	users.Use(middlewares.ValidateAuthToken())

	users.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "User Handle")
	})
	users.GET("/home", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hellow User")
	})

	// signup api
}
