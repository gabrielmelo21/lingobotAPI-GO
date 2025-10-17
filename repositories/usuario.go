package repositories

import (
	"context"
	"fmt"
	"lingobotAPI-GO/config"
	"lingobotAPI-GO/models"
	"lingobotAPI-GO/utils"
	"log"
)

// InsertUsuario insere um novo usuário em todas as 6 tabelas
func InsertUsuario(
	usuario *models.Usuario,
	economia *models.UsuarioEconomia,
	progresso *models.UsuarioProgresso,
	social *models.UsuarioSocial,
	conteudo *models.UsuarioConteudo,
) error {
	ctx := context.Background()

	// Inicia transação
	tx, err := config.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback(ctx)

	// 1. Insere na tabela usuario (retorna ID e created_at)
	var usuarioID int
	queryUsuario := `
		INSERT INTO usuario (nome, sobrenome, email, password, gender, data_nascimento)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	err = tx.QueryRow(ctx, queryUsuario,
		usuario.Nome,
		usuario.Sobrenome,
		usuario.Email,
		usuario.Password,
		usuario.Gender,
		usuario.DataNascimento,
	).Scan(&usuarioID, &usuario.CreatedAt)

	if err != nil {
		return fmt.Errorf("erro ao inserir usuario: %v", err)
	}

	// 2. Insere na tabela usuario_seguranca
	querySeguranca := `
		INSERT INTO usuario_seguranca (usuario_id, otp_code, otp_ativo)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, querySeguranca, usuarioID, nil, false)
	if err != nil {
		return fmt.Errorf("erro ao inserir usuario_seguranca: %v", err)
	}

	// 3. Insere na tabela usuario_economia
	queryEconomia := `
		INSERT INTO usuario_economia (usuario_id, tokens, gemas, battery, plano)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, queryEconomia,
		usuarioID,
		economia.Tokens,
		economia.Gemas,
		economia.Battery,
		economia.Plano,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir usuario_economia: %v", err)
	}

	// 4. Insere na tabela usuario_progresso
	queryProgresso := `
		INSERT INTO usuario_progresso (
			usuario_id, lingo_exp, level, listening, writing, reading,
			speaking, ranking, difficulty, learning
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = tx.Exec(ctx, queryProgresso,
		usuarioID,
		progresso.LingoEXP,
		progresso.Level,
		progresso.Listening,
		progresso.Writing,
		progresso.Reading,
		progresso.Speaking,
		progresso.Ranking,
		progresso.Difficulty,
		progresso.Learning,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir usuario_progresso: %v", err)
	}

	// 5. Insere na tabela usuario_social
	querySocial := `
		INSERT INTO usuario_social (usuario_id, referal_code, invited_by)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, querySocial,
		usuarioID,
		social.ReferalCode,
		social.InvitedBy,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir usuario_social: %v", err)
	}

	// 6. Insere na tabela usuario_conteudo usando Sonic
	itemsJSON, err := utils.Marshal(conteudo.Items)
	if err != nil {
		return fmt.Errorf("erro ao serializar items: %v", err)
	}

	dailyJSON, err := utils.Marshal(conteudo.DailyMissions)
	if err != nil {
		return fmt.Errorf("erro ao serializar daily_missions: %v", err)
	}

	achievementsJSON, err := utils.Marshal(conteudo.Achievements)
	if err != nil {
		return fmt.Errorf("erro ao serializar achievements: %v", err)
	}

	queryConteudo := `
		INSERT INTO usuario_conteudo (usuario_id, items, daily_missions, achievements)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(ctx, queryConteudo,
		usuarioID,
		itemsJSON,
		dailyJSON,
		achievementsJSON,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir usuario_conteudo: %v", err)
	}

	// Commit da transação
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao commitar transação: %v", err)
	}

	// Atualiza o ID do usuário no objeto original
	usuario.ID = usuarioID

	return nil
}

// GetAllUsuarios retorna todos os usuários com dados públicos
func GetAllUsuarios() ([]models.UsuarioResponse, error) {
	ctx := context.Background()

	query := `
		SELECT 
			u.id, u.nome, u.sobrenome, u.gender, u.data_nascimento, u.created_at,
			ue.id, ue.usuario_id, ue.tokens, ue.gemas, ue.battery, ue.plano, ue.updated_at,
			up.id, up.usuario_id, up.lingo_exp, up.level, up.listening, up.writing,
			up.reading, up.speaking, up.ranking, up.difficulty, up.learning, up.updated_at,
			us.id, us.usuario_id, us.referal_code, us.invited_by, us.updated_at,
			uc.id, uc.usuario_id, uc.items, uc.daily_missions, uc.achievements, uc.updated_at
		FROM usuario u
		LEFT JOIN usuario_economia ue ON u.id = ue.usuario_id
		LEFT JOIN usuario_progresso up ON u.id = up.usuario_id
		LEFT JOIN usuario_social us ON u.id = us.usuario_id
		LEFT JOIN usuario_conteudo uc ON u.id = uc.usuario_id
	`

	rows, err := config.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []models.UsuarioResponse

	for rows.Next() {
		var u models.UsuarioResponse
		var economia models.UsuarioEconomia
		var progresso models.UsuarioProgresso
		var social models.UsuarioSocial
		var conteudo models.UsuarioConteudo
		var itemsJSON, dailyJSON, achievementsJSON []byte

		err := rows.Scan(
			&u.ID, &u.Nome, &u.Sobrenome, &u.Gender, &u.DataNascimento, &u.CreatedAt,
			&economia.ID, &economia.UsuarioID, &economia.Tokens, &economia.Gemas,
			&economia.Battery, &economia.Plano, &economia.UpdatedAt,
			&progresso.ID, &progresso.UsuarioID, &progresso.LingoEXP, &progresso.Level,
			&progresso.Listening, &progresso.Writing, &progresso.Reading, &progresso.Speaking,
			&progresso.Ranking, &progresso.Difficulty, &progresso.Learning, &progresso.UpdatedAt,
			&social.ID, &social.UsuarioID, &social.ReferalCode, &social.InvitedBy, &social.UpdatedAt,
			&conteudo.ID, &conteudo.UsuarioID, &itemsJSON, &dailyJSON, &achievementsJSON, &conteudo.UpdatedAt,
		)
		if err != nil {
			log.Printf("Erro ao escanear usuário: %v", err)
			continue
		}

		// Deserializa JSONs usando Sonic
		if err := utils.Unmarshal(itemsJSON, &conteudo.Items); err != nil {
			log.Printf("Erro ao deserializar items: %v", err)
		}
		if err := utils.Unmarshal(dailyJSON, &conteudo.DailyMissions); err != nil {
			log.Printf("Erro ao deserializar daily_missions: %v", err)
		}
		if err := utils.Unmarshal(achievementsJSON, &conteudo.Achievements); err != nil {
			log.Printf("Erro ao deserializar achievements: %v", err)
		}

		u.Economia = economia
		u.Progresso = progresso
		u.Social = &social
		u.Conteudo = &conteudo

		usuarios = append(usuarios, u)
	}

	return usuarios, nil
}

// GetUsuarioByEmail busca usuário completo por email
func GetUsuarioByEmail(email string) (*models.UsuarioCompleto, error) {
	ctx := context.Background()

	query := `
		SELECT 
			u.id, u.nome, u.sobrenome, u.email, u.password, u.gender, u.data_nascimento, u.created_at,
			useg.id, useg.usuario_id, useg.otp_code, useg.otp_ativo, useg.updated_at,
			ue.id, ue.usuario_id, ue.tokens, ue.gemas, ue.battery, ue.plano, ue.updated_at,
			up.id, up.usuario_id, up.lingo_exp, up.level, up.listening, up.writing,
			up.reading, up.speaking, up.ranking, up.difficulty, up.learning, up.updated_at,
			us.id, us.usuario_id, us.referal_code, us.invited_by, us.updated_at,
			uc.id, uc.usuario_id, uc.items, uc.daily_missions, uc.achievements, uc.updated_at
		FROM usuario u
		LEFT JOIN usuario_seguranca useg ON u.id = useg.usuario_id
		LEFT JOIN usuario_economia ue ON u.id = ue.usuario_id
		LEFT JOIN usuario_progresso up ON u.id = up.usuario_id
		LEFT JOIN usuario_social us ON u.id = us.usuario_id
		LEFT JOIN usuario_conteudo uc ON u.id = uc.usuario_id
		WHERE u.email = $1
	`

	var uc models.UsuarioCompleto
	var seguranca models.UsuarioSeguranca
	var economia models.UsuarioEconomia
	var progresso models.UsuarioProgresso
	var social models.UsuarioSocial
	var conteudo models.UsuarioConteudo
	var itemsJSON, dailyJSON, achievementsJSON []byte

	err := config.DB.QueryRow(ctx, query, email).Scan(
		&uc.Usuario.ID, &uc.Usuario.Nome, &uc.Usuario.Sobrenome, &uc.Usuario.Email,
		&uc.Usuario.Password, &uc.Usuario.Gender, &uc.Usuario.DataNascimento, &uc.Usuario.CreatedAt,
		&seguranca.ID, &seguranca.UsuarioID, &seguranca.OTPCode, &seguranca.OTPAtivo, &seguranca.UpdatedAt,
		&economia.ID, &economia.UsuarioID, &economia.Tokens, &economia.Gemas,
		&economia.Battery, &economia.Plano, &economia.UpdatedAt,
		&progresso.ID, &progresso.UsuarioID, &progresso.LingoEXP, &progresso.Level,
		&progresso.Listening, &progresso.Writing, &progresso.Reading, &progresso.Speaking,
		&progresso.Ranking, &progresso.Difficulty, &progresso.Learning, &progresso.UpdatedAt,
		&social.ID, &social.UsuarioID, &social.ReferalCode, &social.InvitedBy, &social.UpdatedAt,
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

	uc.Seguranca = &seguranca
	uc.Economia = economia
	uc.Progresso = progresso
	uc.Social = &social
	uc.Conteudo = &conteudo

	return &uc, nil
}

// GetUsuarioByID busca usuário completo por ID
func GetUsuarioByID(id int) (*models.UsuarioCompleto, error) {
	ctx := context.Background()

	query := `
		SELECT 
			u.id, u.nome, u.sobrenome, u.email, u.password, u.gender, u.data_nascimento, u.created_at,
			useg.id, useg.usuario_id, useg.otp_code, useg.otp_ativo, useg.updated_at,
			ue.id, ue.usuario_id, ue.tokens, ue.gemas, ue.battery, ue.plano, ue.updated_at,
			up.id, up.usuario_id, up.lingo_exp, up.level, up.listening, up.writing,
			up.reading, up.speaking, up.ranking, up.difficulty, up.learning, up.updated_at,
			us.id, us.usuario_id, us.referal_code, us.invited_by, us.updated_at,
			uc.id, uc.usuario_id, uc.items, uc.daily_missions, uc.achievements, uc.updated_at
		FROM usuario u
		LEFT JOIN usuario_seguranca useg ON u.id = useg.usuario_id
		LEFT JOIN usuario_economia ue ON u.id = ue.usuario_id
		LEFT JOIN usuario_progresso up ON u.id = up.usuario_id
		LEFT JOIN usuario_social us ON u.id = us.usuario_id
		LEFT JOIN usuario_conteudo uc ON u.id = uc.usuario_id
		WHERE u.id = $1
	`

	var uc models.UsuarioCompleto
	var seguranca models.UsuarioSeguranca
	var economia models.UsuarioEconomia
	var progresso models.UsuarioProgresso
	var social models.UsuarioSocial
	var conteudo models.UsuarioConteudo
	var itemsJSON, dailyJSON, achievementsJSON []byte

	err := config.DB.QueryRow(ctx, query, id).Scan(
		&uc.Usuario.ID, &uc.Usuario.Nome, &uc.Usuario.Sobrenome, &uc.Usuario.Email,
		&uc.Usuario.Password, &uc.Usuario.Gender, &uc.Usuario.DataNascimento, &uc.Usuario.CreatedAt,
		&seguranca.ID, &seguranca.UsuarioID, &seguranca.OTPCode, &seguranca.OTPAtivo, &seguranca.UpdatedAt,
		&economia.ID, &economia.UsuarioID, &economia.Tokens, &economia.Gemas,
		&economia.Battery, &economia.Plano, &economia.UpdatedAt,
		&progresso.ID, &progresso.UsuarioID, &progresso.LingoEXP, &progresso.Level,
		&progresso.Listening, &progresso.Writing, &progresso.Reading, &progresso.Speaking,
		&progresso.Ranking, &progresso.Difficulty, &progresso.Learning, &progresso.UpdatedAt,
		&social.ID, &social.UsuarioID, &social.ReferalCode, &social.InvitedBy, &social.UpdatedAt,
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

	uc.Seguranca = &seguranca
	uc.Economia = economia
	uc.Progresso = progresso
	uc.Social = &social
	uc.Conteudo = &conteudo

	return &uc, nil
}

// UpdateUsuario atualiza dados do usuário (PLACEHOLDER - implementar depois)
func UpdateUsuario(u *models.Usuario) error {
	// TODO: Implementar update nas tabelas necessárias
	return nil
}

// UpdateUsuarioCompleto atualiza todas as tabelas do usuário
func UpdateUsuarioCompleto(uc *models.UsuarioCompleto) error {
	ctx := context.Background()

	// Inicia transação
	tx, err := config.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback(ctx)

	// 1. Atualiza tabela usuario
	queryUsuario := `
		UPDATE usuario SET
			nome = $1, sobrenome = $2, email = $3,
			gender = $4, data_nascimento = $5
		WHERE id = $6
	`
	_, err = tx.Exec(ctx, queryUsuario,
		uc.Usuario.Nome,
		uc.Usuario.Sobrenome,
		uc.Usuario.Email,
		uc.Usuario.Gender,
		uc.Usuario.DataNascimento,
		uc.Usuario.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar usuario: %v", err)
	}

	// 2. Atualiza tabela usuario_economia
	queryEconomia := `
		UPDATE usuario_economia SET
			tokens = $1, gemas = $2, battery = $3, plano = $4, updated_at = CURRENT_TIMESTAMP
		WHERE usuario_id = $5
	`
	_, err = tx.Exec(ctx, queryEconomia,
		uc.Economia.Tokens,
		uc.Economia.Gemas,
		uc.Economia.Battery,
		uc.Economia.Plano,
		uc.Usuario.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar usuario_economia: %v", err)
	}

	// 3. Atualiza tabela usuario_progresso
	queryProgresso := `
		UPDATE usuario_progresso SET
			lingo_exp = $1, level = $2, listening = $3, writing = $4,
			reading = $5, speaking = $6, ranking = $7,
			difficulty = $8, learning = $9, updated_at = CURRENT_TIMESTAMP
		WHERE usuario_id = $10
	`
	_, err = tx.Exec(ctx, queryProgresso,
		uc.Progresso.LingoEXP,
		uc.Progresso.Level,
		uc.Progresso.Listening,
		uc.Progresso.Writing,
		uc.Progresso.Reading,
		uc.Progresso.Speaking,
		uc.Progresso.Ranking,
		uc.Progresso.Difficulty,
		uc.Progresso.Learning,
		uc.Usuario.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar usuario_progresso: %v", err)
	}

	// 4. Atualiza tabela usuario_conteudo usando Sonic
	itemsJSON, err := utils.Marshal(uc.Conteudo.Items)
	if err != nil {
		return fmt.Errorf("erro ao serializar items: %v", err)
	}

	dailyJSON, err := utils.Marshal(uc.Conteudo.DailyMissions)
	if err != nil {
		return fmt.Errorf("erro ao serializar daily_missions: %v", err)
	}

	achievementsJSON, err := utils.Marshal(uc.Conteudo.Achievements)
	if err != nil {
		return fmt.Errorf("erro ao serializar achievements: %v", err)
	}

	queryConteudo := `
		UPDATE usuario_conteudo SET
			items = $1, daily_missions = $2, achievements = $3, updated_at = CURRENT_TIMESTAMP
		WHERE usuario_id = $4
	`
	_, err = tx.Exec(ctx, queryConteudo,
		itemsJSON,
		dailyJSON,
		achievementsJSON,
		uc.Usuario.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar usuario_conteudo: %v", err)
	}

	// Commit da transação
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao commitar transação: %v", err)
	}

	return nil
}
