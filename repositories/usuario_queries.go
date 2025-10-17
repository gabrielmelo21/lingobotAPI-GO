package repositories

import (
	"context"
	"fmt"
	"lingobotAPI-GO/config"
	"lingobotAPI-GO/models"
	"lingobotAPI-GO/utils"
)

// GetUsuarioProfile retorna dados básicos do perfil (tabela usuario)
func GetUsuarioProfile(usuarioID int) (*models.Usuario, error) {
	ctx := context.Background()

	query := `
		SELECT id, nome, sobrenome, email, gender, data_nascimento, created_at
		FROM usuario
		WHERE id = $1
	`

	var u models.Usuario
	err := config.DB.QueryRow(ctx, query, usuarioID).Scan(
		&u.ID, &u.Nome, &u.Sobrenome, &u.Email,
		&u.Gender, &u.DataNascimento, &u.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUsuarioContent retorna economia, progresso e conteúdo (sem dados sensíveis)
func GetUsuarioContent(usuarioID int) (*models.UsuarioContentResponse, error) {
	ctx := context.Background()

	query := `
		SELECT 
			ue.id, ue.usuario_id, ue.tokens, ue.gemas, ue.battery, ue.plano, ue.updated_at,
			up.id, up.usuario_id, up.lingo_exp, up.level, up.listening, up.writing,
			up.reading, up.speaking, up.ranking, up.difficulty, up.learning, up.updated_at,
			uc.id, uc.usuario_id, uc.items, uc.daily_missions, uc.achievements, uc.updated_at
		FROM usuario_economia ue
		LEFT JOIN usuario_progresso up ON ue.usuario_id = up.usuario_id
		LEFT JOIN usuario_conteudo uc ON ue.usuario_id = uc.usuario_id
		WHERE ue.usuario_id = $1
	`

	var economia models.UsuarioEconomia
	var progresso models.UsuarioProgresso
	var conteudo models.UsuarioConteudo
	var itemsJSON, dailyJSON, achievementsJSON []byte

	err := config.DB.QueryRow(ctx, query, usuarioID).Scan(
		&economia.ID, &economia.UsuarioID, &economia.Tokens, &economia.Gemas,
		&economia.Battery, &economia.Plano, &economia.UpdatedAt,
		&progresso.ID, &progresso.UsuarioID, &progresso.LingoEXP, &progresso.Level,
		&progresso.Listening, &progresso.Writing, &progresso.Reading, &progresso.Speaking,
		&progresso.Ranking, &progresso.Difficulty, &progresso.Learning, &progresso.UpdatedAt,
		&conteudo.ID, &conteudo.UsuarioID, &itemsJSON, &dailyJSON, &achievementsJSON, &conteudo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Deserializa JSONs usando Sonic
	if err := utils.Unmarshal(itemsJSON, &conteudo.Items); err != nil {
		return nil, fmt.Errorf("erro ao deserializar items: %v", err)
	}
	if err := utils.Unmarshal(dailyJSON, &conteudo.DailyMissions); err != nil {
		return nil, fmt.Errorf("erro ao deserializar daily_missions: %v", err)
	}
	if err := utils.Unmarshal(achievementsJSON, &conteudo.Achievements); err != nil {
		return nil, fmt.Errorf("erro ao deserializar achievements: %v", err)
	}

	return &models.UsuarioContentResponse{
		Economia:  economia,
		Progresso: progresso,
		Conteudo:  conteudo,
	}, nil
}

// GetUsuarioSocial retorna dados sociais (referal_code, invited_by)
func GetUsuarioSocial(usuarioID int) (*models.UsuarioSocial, error) {
	ctx := context.Background()

	query := `
		SELECT id, usuario_id, referal_code, invited_by, updated_at
		FROM usuario_social
		WHERE usuario_id = $1
	`

	var social models.UsuarioSocial
	err := config.DB.QueryRow(ctx, query, usuarioID).Scan(
		&social.ID, &social.UsuarioID, &social.ReferalCode,
		&social.InvitedBy, &social.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &social, nil
}

// GetUsuarioSecurity retorna dados de segurança (OTP) - ADMIN ONLY
func GetUsuarioSecurity(usuarioID int) (*models.UsuarioSeguranca, error) {
	ctx := context.Background()

	query := `
		SELECT id, usuario_id, otp_code, otp_ativo, updated_at
		FROM usuario_seguranca
		WHERE usuario_id = $1
	`

	var seg models.UsuarioSeguranca
	err := config.DB.QueryRow(ctx, query, usuarioID).Scan(
		&seg.ID, &seg.UsuarioID, &seg.OTPCode,
		&seg.OTPAtivo, &seg.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &seg, nil
}
