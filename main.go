package main

import (
	"lingobotAPI-GO/config"
	"lingobotAPI-GO/routes"
	"os"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // porta padrão local
	}

	err := router.Run(":" + port)
	if err != nil {
		return
	}

}
