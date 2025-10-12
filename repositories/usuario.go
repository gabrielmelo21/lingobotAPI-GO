package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	_ "fmt"
	"lingobotAPI-GO/config"
	"lingobotAPI-GO/models"
	"log"
)

func InsertUsuario(u *models.Usuario) error {
	ctx := context.Background()

	query := `INSERT INTO usuario (
		nome, sobrenome, email, password, "OTP_code", "LingoEXP", "Level",
		gender, data_nascimento, tokens, plano, created_at, referal_code,
		invited_by, ranking, listening, writing, reading, speaking, gemas,
		items, "dailyMissions", achievements, difficulty, battery, learning
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
		$14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
	) RETURNING id`

	itemsJSON, err := json.Marshal(u.Items)
	if err != nil {
		return fmt.Errorf("erro ao converter Items: %w", err)
	}
	dailyJSON, err := json.Marshal(u.DailyMissions)
	if err != nil {
		return fmt.Errorf("erro ao converter DailyMissions: %w", err)
	}
	achievementsJSON, err := json.Marshal(u.Achievements)
	if err != nil {
		return fmt.Errorf("erro ao converter Achievements: %w", err)
	}

	err = config.DB.QueryRow(ctx, query,
		u.Nome, u.Sobrenome, u.Email, u.Password, u.OTPCode, u.LingoEXP,
		u.Level, u.Gender, u.DataNascimento, u.Tokens, u.Plano, u.CreatedAt,
		u.ReferalCode, u.InvitedBy, u.Ranking, u.Listening, u.Writing,
		u.Reading, u.Speaking, u.Gemas, itemsJSON, dailyJSON,
		achievementsJSON, u.Difficulty, u.Battery, u.Learning,
	).Scan(&u.ID)

	if err != nil {
		return fmt.Errorf("erro ao inserir usuário: %w", err)
	}

	return nil
}

func GetAllUsuarios() ([]models.UsuarioResponse, error) {
	ctx := context.Background()
	rows, err := config.DB.Query(ctx, `SELECT
		id, nome, sobrenome, "LingoEXP", "Level",
		gender, data_nascimento, tokens, plano, created_at, referal_code,
		invited_by, ranking, listening, writing, reading, speaking, gemas,
		items, "dailyMissions", achievements, difficulty, battery, learning
	FROM usuario`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []models.UsuarioResponse
	for rows.Next() {
		var u models.UsuarioResponse
		var itemsJSON, dailyJSON, achievementsJSON []byte
		var dataNascimentoStr *string

		err := rows.Scan(
			&u.ID, &u.Nome, &u.Sobrenome, &u.LingoEXP, &u.Level,
			&u.Gender, &dataNascimentoStr, &u.Tokens,
			&u.Plano, &u.CreatedAt, &u.ReferalCode, &u.InvitedBy, &u.Ranking,
			&u.Listening, &u.Writing, &u.Reading, &u.Speaking, &u.Gemas,
			&itemsJSON, &dailyJSON, &achievementsJSON, &u.Difficulty,
			&u.Battery, &u.Learning,
		)
		if err != nil {
			log.Printf("Erro ao escanear usuário: %v", err)
			continue
		}

		// Deserializa os JSONs
		err = json.Unmarshal(itemsJSON, &u.Items)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dailyJSON, &u.DailyMissions)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(achievementsJSON, &u.Achievements)
		if err != nil {
			return nil, err
		}

		usuarios = append(usuarios, u)
	}
	return usuarios, nil
}

// GetUsuarioByEmail busca usuário por email
func GetUsuarioByEmail(email string) (*models.Usuario, error) {
	ctx := context.Background()
	var u models.Usuario
	var dataNascimentoStr *string

	query := `SELECT
		id, nome, sobrenome, email, password, "OTP_code", "LingoEXP", "Level",
		gender, data_nascimento, tokens, plano, created_at, referal_code,
		invited_by, ranking, listening, writing, reading, speaking, gemas,
		items, "dailyMissions", achievements, difficulty, battery, learning
	FROM usuario WHERE email = $1`

	var itemsJSON, dailyJSON, achievementsJSON []byte

	err := config.DB.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Nome, &u.Sobrenome, &u.Email, &u.Password, &u.OTPCode,
		&u.LingoEXP, &u.Level, &u.Gender, &dataNascimentoStr, &u.Tokens,
		&u.Plano, &u.CreatedAt, &u.ReferalCode, &u.InvitedBy, &u.Ranking,
		&u.Listening, &u.Writing, &u.Reading, &u.Speaking, &u.Gemas,
		&itemsJSON, &dailyJSON, &achievementsJSON, &u.Difficulty,
		&u.Battery, &u.Learning,
	)

	if err != nil {
		return nil, err
	}

	// Deserializa os JSONs
	err = json.Unmarshal(itemsJSON, &u.Items)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dailyJSON, &u.DailyMissions)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(achievementsJSON, &u.Achievements)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
