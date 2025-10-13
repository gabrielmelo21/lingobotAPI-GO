package routes

import (
	"lingobotAPI-GO/controllers"
	"lingobotAPI-GO/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	// Rotas públicas (sem autenticação)
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
