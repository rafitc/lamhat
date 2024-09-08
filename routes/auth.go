package routes

import (
	"lamhat/model"
	"lamhat/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddAuthRoutes(rg *gin.RouterGroup) {

	Sugar_logger.Debug("Request received in Auth group")
	auth := rg.Group("/")

	// signup api
	auth.POST("/signup", func(ctx *gin.Context) {

		var body model.SignupBody
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var result model.Response = service.SignUpService(ctx, body)
		ctx.JSON(result.Code, result)
	})

	// login API
	auth.POST("/login", func(ctx *gin.Context) {
		var body model.LoginBody
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var result model.Response = service.LoginService(ctx, body)
		ctx.JSON(result.Code, result)
	})
}
