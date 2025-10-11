package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	router := gin.Default()

	// Configuração CORS específica para apps com Capacitor
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"capacitor://localhost",
			"http://localhost",      // para testes web
			"http://localhost:4200", // para Angular rodando local
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	errToLoadEnv := godotenv.Load()
	if errToLoadEnv != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	// Lê variáveis do .env
	port := os.Getenv("PORT")
	greeting := os.Getenv("GREETING")

	fmt.Println("Porta:", port)
	fmt.Println("Mensagem:", greeting)

	// 1️⃣ JSON de resposta
	router.GET("/json", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"mensagem": "Aqui temos um JSON com gin",
		})
	})

	// 2️⃣ Texto puro
	router.GET("/text", func(c *gin.Context) {
		c.String(200, "Aqui temos um texto puro")
	})

	// 3️⃣ Query string
	// Ex: /query?nome=Gabriel
	router.GET("/query", func(c *gin.Context) {
		nome := c.Query("nome")
		if nome == "" {
			nome = "Visitante"
		}
		c.String(200, "Olá, %s!", nome)
	})

	// 4️⃣ Parâmetro de rota
	// Ex: /param/Gabriel
	router.GET("/param/:nome", func(c *gin.Context) {
		nome := c.Param("nome") // sem os dois pontos aqui
		c.String(200, "Olá, %s! Esse é um parâmetro de rota.", nome)
	})

	// 5️⃣ Exemplo de POST recebendo JSON
	router.POST("/postjson", func(c *gin.Context) {
		var json struct {
			Mensagem string `json:"mensagem"`
		}
		if err := c.BindJSON(&json); err != nil {
			c.JSON(400, gin.H{"erro": "JSON inválido"})
			return
		}
		c.JSON(200, gin.H{
			"mensagem_recebida": json.Mensagem,
		})
	})

	errToRun := router.Run(":8081")
	if errToRun != nil {
		return
	}

}
