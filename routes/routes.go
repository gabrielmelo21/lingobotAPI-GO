package routes

import (
	"controllers/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/hello", controllers.Hello)
}
