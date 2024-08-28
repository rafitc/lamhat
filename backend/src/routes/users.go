package routes

import (
	"backend/src/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Sugar_logger = core.Sugar

func AddUserRoutes(rg *gin.RouterGroup) {

	Sugar_logger.Debug("Request received in Users group")
	users := rg.Group("/")

	users.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "users")
	})
	users.GET("/welcome", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hellow User")
	})
}
