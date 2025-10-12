package routes

import (
	"github.com/gin-gonic/gin"
	"lingobotAPI-GO/controllers"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/usuarios", controllers.GetUsuarios)
	router.POST("/usuarios", controllers.CriarUsuario)
	router.POST("/login", controllers.Login)
}
