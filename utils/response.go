package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SonicJSON é uma função auxiliar para enviar respostas JSON usando o Sonic.
// Ela substitui a necessidade de chamar c.JSON() diretamente.
func SonicJSON(c *gin.Context, code int, obj interface{}) {
	bytes, err := Marshal(obj) // Usando a função Marshal do utils/json.go
	if err != nil {
		// Idealmente, logar o erro em um sistema de logs
		c.String(http.StatusInternalServerError, "Internal Server Error: "+err.Error())
		return
	}
	c.Data(code, "application/json; charset=utf-8", bytes)
}
