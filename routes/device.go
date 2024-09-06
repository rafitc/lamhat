package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddDeviceRoutes(rg *gin.RouterGroup) {

	Sugar_logger.Debug("Request received in Users group")
	users := rg.Group("/")

	users.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "device")
	})
	users.GET("/welcome", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hellow Device welcome")
	})
}
