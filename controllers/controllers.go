package controllers

import (
	"context"
	"lingobotAPI-GO/config"
	"lingobotAPI-GO/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsuarios(c *gin.Context) {
	ctx := context.Background()

	rows, err := config.DB.Query(ctx, `
    SELECT 
        id, 
        nome, 
        sobrenome, 
        email, 
        password, 
        "OTP_code", 
        "LingoEXP", 
        "Level", 
        gender, 
        data_nascimento,
        tokens, 
        plano, 
        created_at, 
        referal_code, 
        invited_by, 
        ranking, 
        listening, 
        writing,
        reading, 
        speaking, 
        gemas, 
        items, 
        "dailyMissions", 
        achievements, 
        difficulty, 
        battery, 
        learning
    FROM usuario
`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao consultar usuários", "detalhe": err.Error()})
		return
	}
	defer rows.Close()

	var usuarios []models.Usuario

	for rows.Next() {
		var u models.Usuario
		err := rows.Scan(
			&u.ID, &u.Nome, &u.Sobrenome, &u.Email, &u.Password, &u.OTPCode, &u.LingoEXP, &u.Level,
			&u.Gender, &u.DataNascimento, &u.Tokens, &u.Plano, &u.CreatedAt, &u.ReferalCode,
			&u.InvitedBy, &u.Ranking, &u.Listening, &u.Writing, &u.Reading, &u.Speaking,
			&u.Gemas, &u.Items, &u.DailyMissions, &u.Achievements, &u.Difficulty, &u.Battery, &u.Learning,
		)
		if err != nil {
			log.Printf("Erro ao consultar usuários: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
			return
		}

		usuarios = append(usuarios, u)
	}

	c.JSON(http.StatusOK, usuarios)
}
