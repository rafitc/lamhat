package routes

import (
	"backend/src/core"
	"backend/src/model"
	"backend/src/service"
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

	// signup api
	users.POST("/signup", func(ctx *gin.Context) {

		var body model.SignupBody
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var result model.Response = service.SingUpService(ctx, body)
		ctx.JSON(result.Code, result)
	})
}
