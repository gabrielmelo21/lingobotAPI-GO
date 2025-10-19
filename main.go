package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"lingobotAPI-GO/config"
	"lingobotAPI-GO/routes"
	"lingobotAPI-GO/utils"
)

func main() {
	router := gin.Default()

	// Substitui o binder JSON padrão do Gin pelo nosso binder customizado (Sonic)
	binding.JSON = utils.NewJsonBinding()

	// Configuração do CORS (vem do config/cors.go)
	config.SetupCORS(router)

	// Inicializar o banco de dados (vem do config/database.go)
	config.ConnectDatabase()

	// Registrar as rotas (vem de routes/routes.go)
	routes.RegisterRoutes(router)

	port := "8100"

	err := router.Run(":" + port)
	if err != nil {
		return
	}

}