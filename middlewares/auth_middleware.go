package middlewares

import (
	"lingobotAPI-GO/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware valida o JWT antes de permitir acesso à rota
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pega o header Authorization
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"erro": "Token não fornecido"})
			c.Abort()
			return
		}

		// Verifica se começa com "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"erro": "Formato de token inválido"})
			c.Abort()
			return
		}

		// Extrai o token (remove "Bearer ")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Valida o token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"erro": "Token inválido ou expirado"})
			c.Abort()
			return
		}

		// Armazena os claims no contexto para uso posterior
		c.Set("user_id", claims.Sub)
		c.Set("claims", claims)

		// Continua para o próximo handler
		c.Next()
	}
}
