package main

import (
	"lingobotAPI-GO/config"
	"lingobotAPI-GO/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Configuração do CORS (vem do config/cors.go)
	config.SetupCORS(router)

	// Inicializar o banco de dados (vem do config/database.go)
	config.ConnectDatabase()

	// Registrar as rotas (vem de routes/routes.go)
	routes.RegisterRoutes(router)

	// Rodar o servidor
	router.Run(":8100")
}
