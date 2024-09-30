package routes

import (
	"lamhat/middlewares"
	"lamhat/model"
	"lamhat/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddGalleryRoutes(rg *gin.RouterGroup) {

	Sugar_logger.Debug("Request received in Gallery group")
	gallery := rg.Group("/")
	gallery.Use(middlewares.ValidateAuthToken())

	gallery.GET("/get/", func(ctx *gin.Context) {
		gallery_id, err := strconv.Atoi(ctx.Query("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user_id := ctx.GetInt("userId")
		var result model.Response = service.FetchGallery(ctx, gallery_id, user_id)
		ctx.JSON(result.Code, result)
	})

	gallery.POST("/create-gallery", func(ctx *gin.Context) {
		var body model.CreateGallery
		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var result model.Response = service.CreateGallery(ctx, body)
		ctx.JSON(result.Code, result)
	})

	gallery.POST("/upload/", func(ctx *gin.Context) {
		gallery_id, err := strconv.Atoi(ctx.Query("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Handle Files
		user_id := ctx.GetInt("userId")
		var result model.Response = service.UploadIntoGallery(ctx, gallery_id, user_id)
		ctx.JSON(result.Code, result)
	})

	gallery.POST("/change-status", func(ctx *gin.Context) {
		var body model.PublishGallery

		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user_id := ctx.GetInt("userId")
		var result model.Response = service.PublishGallery(ctx, user_id, body)
		ctx.JSON(result.Code, result)
	})
}
