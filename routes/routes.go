package routes

import (
	"lingobotAPI-GO/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/usuarios", controllers.GetUsuarios)
}
