package main

import (
	"lamhat/core"
	"lamhat/repository"
	"lamhat/routes"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()
var sugar = core.Sugar

func main() {
	sugar.Info("Creating DB pool")
	repository.DbPoolMain()

	getRoutes()

	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	// })
	sugar.Info("Starting server...")

	router.Run()
}

// getRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func getRoutes() {
	sugar.Debug("Registering Routes")

	auth := router.Group("/auth")
	routes.AddAuthRoutes(auth)

	users := router.Group("/users")
	routes.AddUserRoutes(users)

	devices := router.Group("/device")
	routes.AddDeviceRoutes(devices)

	gallery := router.Group("/gallery")
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	routes.AddGalleryRoutes(gallery)
}
