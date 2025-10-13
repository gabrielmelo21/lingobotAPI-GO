package routes

import (
	"lingobotAPI-GO/controllers"
	"lingobotAPI-GO/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	// Rotas públicas (sem autenticação)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"uptime":  time.Since(time.Now()).String(),
			"version": "1.0.0",
		})
	})

	router.POST("/usuarios", controllers.CriarUsuario)
	router.POST("/login", controllers.Login)

	// Rotas protegidas (com autenticação)
	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/update-user-data", controllers.UpdateUserData)

		// precisa de ma proteção plus ADM apenas
		//protected.GET("/usuarios", controllers.GetUsuarios)

		// IA - Todas as rotas protegidas
		protected.POST("/ai/gemini", controllers.AIGemini)
		protected.POST("/ai/cohere", controllers.AICohere)
		protected.POST("/ai/mistral", controllers.AIMistral)
		protected.POST("/ai/groq", controllers.AIGroq)
		protected.POST("/ai/openrouter", controllers.AIOpenRouter)
		protected.POST("/ai/benchmark", controllers.AIBenchmark)

		// Mídia - TTS e Transcrição
		protected.POST("/tts", controllers.TTS)
		protected.POST("/transcribe", controllers.TranscribeAudio)

	}
}
