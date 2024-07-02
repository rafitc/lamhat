package main

import (
	core "backend/src/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	sugar := core.Sugar

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	sugar.Info("Starting server...")

	r.Run()
}
