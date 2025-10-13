package controllers

import (
	"lingobotAPI-GO/models"
	"lingobotAPI-GO/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AIGemini endpoint principal com fallback
func AIGemini(c *gin.Context) {
	var req models.AIRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text input is required"})
		return
	}

	response, err := services.CallAIWithFallback(req.Text, req.Mistral, req.Cohere, req.Groq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retorna texto puro
	c.String(http.StatusOK, response)
}

// AICohere endpoint específico para Cohere
func AICohere(c *gin.Context) {
	var req models.AIRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text input is required"})
		return
	}

	response, err := services.CallCohere(req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, response)
}

// AIMistral endpoint específico para Mistral
func AIMistral(c *gin.Context) {
	var req models.AIRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text input is required"})
		return
	}

	response, err := services.CallMistral(req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, response)
}

// AIGroq endpoint específico para Groq
func AIGroq(c *gin.Context) {
	var req models.AIRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text input is required"})
		return
	}

	response, err := services.CallGroq(req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, response)
}

// AIOpenRouter endpoint específico para OpenRouter
func AIOpenRouter(c *gin.Context) {
	var req models.AIRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text input is required"})
		return
	}

	response, err := services.CallOpenRouter(req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, response)
}

// AIBenchmark testa todas as IAs e retorna tempos de resposta
func AIBenchmark(c *gin.Context) {
	var req models.AIRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text input is required"})
		return
	}

	results := make(models.BenchmarkResponse)

	// Função auxiliar para benchmark
	benchmark := func(name string, fn func(string) (string, error)) {
		start := time.Now()
		response, err := fn(req.Text)
		duration := time.Since(start).Seconds()

		if err != nil {
			results[name] = models.AIResponse{
				Response: "",
				Error:    err.Error(),
				Time:     duration,
			}
		} else {
			results[name] = models.AIResponse{
				Response: response,
				Time:     duration,
			}
		}
	}

	// Executa benchmark de cada modelo
	benchmark("Gemini", services.CallGemini)
	benchmark("Mistral", services.CallMistral)
	benchmark("Cohere", services.CallCohere)
	benchmark("Groq", services.CallGroq)
	benchmark("OpenRouter", services.CallOpenRouter)

	c.JSON(http.StatusOK, results)
}
