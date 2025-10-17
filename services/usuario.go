package services

import (
	"errors"
	"lingobotAPI-GO/models"
	"lingobotAPI-GO/repositories"
	"lingobotAPI-GO/utils"
	"time"
)

type CriarUsuarioRequest struct {
	Nome           string  `json:"nome" binding:"required"`
	Sobrenome      string  `json:"sobrenome" binding:"required"`
	Email          string  `json:"email" binding:"required"`
	Password       string  `json:"password" binding:"required"`
	Gender         *string `json:"gender"`
	DataNascimento *string `json:"data_nascimento"`
}

// CriarUsuario cria um novo usuário nas 6 tabelas
func CriarUsuario(req CriarUsuarioRequest) error {
	// Validação de nome e sobrenome
	if !utils.ValidateNome(req.Nome) {
		return errors.New("nome deve conter apenas letras")
	}
	if !utils.ValidateNome(req.Sobrenome) {
		return errors.New("sobrenome deve conter apenas letras")
	}

	// Validação de email
	if !utils.ValidateEmail(req.Email) {
		return errors.New("e-mail inválido")
	}

	// Verifica se o usuário já existe
	existente, err := repositories.GetUsuarioByEmail(req.Email)
	if err == nil && existente != nil {
		return errors.New("e-mail já cadastrado")
	}

	// Hash da senha
	senhaHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("erro ao processar senha")
	}

	// Itens iniciais
	itensIniciais := []map[string]interface{}{
		{
			"itemName":  "OG Ticket",
			"dropRate":  0.01,
			"gemsValue": 20,
			"rarity":    "legendary",
			"itemSrc":   "assets/lingobot/itens/og_ticket.webp",
			"describe":  "OG ticket é para os pioneiros.",
			"quant":     1,
		},
		{
			"itemName":  "Beta Tester Ticket",
			"dropRate":  0.01,
			"gemsValue": 20,
			"rarity":    "legendary",
			"itemSrc":   "assets/lingobot/itens/beta_tester_ticket.webp",
			"describe":  "Ticket dos escolhidos.",
			"quant":     1,
		},
	}

	// Daily Missions iniciais
	dailyMissionsIniciais := map[string]interface{}{
		"writing":        false,
		"reading":        false,
		"listening":      false,
		"speaking":       false,
		"chestWasOpen1":  false,
		"chestWasOpen2":  false,
		"chestWasOpen3":  false,
		"chestWasOpen4":  false,
		"strikes":        0,
		"rewardPerChest": 5,
		"chestsOpenedAt": 0,
		"refreshTimeAt":  0,
	}

	// Achievements iniciais (55 falses)
	achievementsList := make([]bool, 55)
	achievementsIniciais := map[string]interface{}{
		"achievements": achievementsList,
	}

	// 1. Prepara dados da tabela usuario
	usuario := &models.Usuario{
		Nome:           req.Nome,
		Sobrenome:      &req.Sobrenome,
		Email:          req.Email,
		Password:       senhaHash,
		Gender:         req.Gender,
		DataNascimento: req.DataNascimento,
		CreatedAt:      time.Now(),
	}

	// 2. Prepara dados da tabela usuario_economia
	economia := &models.UsuarioEconomia{
		Tokens:  0,
		Gemas:   0,
		Battery: 10,
		Plano:   "free",
	}

	// 3. Prepara dados da tabela usuario_progresso
	progresso := &models.UsuarioProgresso{
		LingoEXP:   0,
		Level:      1,
		Listening:  0,
		Writing:    0,
		Reading:    0,
		Speaking:   0,
		Ranking:    4,
		Difficulty: "medium",
		Learning:   "en",
	}

	// 4. Prepara dados da tabela usuario_social
	social := &models.UsuarioSocial{
		ReferalCode: nil, // Será implementado depois se necessário
		InvitedBy:   nil,
	}

	// 5. Prepara dados da tabela usuario_conteudo
	conteudo := &models.UsuarioConteudo{
		Items:         itensIniciais,
		DailyMissions: dailyMissionsIniciais,
		Achievements:  achievementsIniciais,
	}

	// Insere o usuário em todas as tabelas (transação)
	err = repositories.InsertUsuario(usuario, economia, progresso, social, conteudo)
	if err != nil {
		return err
	}

	return nil
}
