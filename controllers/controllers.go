package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"mensagem": "Olá, Gabriel! 🚀 Estrutura organizada com sucesso",
	})
}
