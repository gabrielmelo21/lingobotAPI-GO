package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func ConnectDatabase() {
	// Carrega o .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Aviso: não foi possível carregar o .env, usando variáveis do ambiente")
	}

	// Lê a variável DATABASE_URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("❌ ERRO: DATABASE_URL não encontrada no .env")
	}

	// Cria o pool de conexão
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("❌ Erro ao analisar DATABASE_URL: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DB, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco: %v", err)
	}

	// Testa a conexão
	err = DB.Ping(ctx)
	if err != nil {
		log.Fatalf("❌ Falha ao testar conexão com o banco: %v", err)
	}

	fmt.Println("✅ Conectado ao banco de dados com sucesso!")
}
