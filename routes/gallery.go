package routes

import (
	"lamhat/middlewares"
	"lamhat/model"
	"lamhat/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddGalleryRoutes(rg *gin.RouterGroup) {

	user := rg.Group("/")
	user.Use(middlewares.ValidateAuthToken())

	user.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Gallery Handle")
	})

	user.POST("/create-gallery", func(ctx *gin.Context) {
		var body model.CreateGallery
		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var result model.Response = service.CreateGallery(ctx, body)
		ctx.JSON(result.Code, result)
	})
}
