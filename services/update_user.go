package services

import (
	"errors"
	_ "lingobotAPI-GO/models"
	"lingobotAPI-GO/repositories"
	"lingobotAPI-GO/utils"
)

type UpdateUserDataRequest struct {
	ID             *int        `json:"id"`
	Sub            *int        `json:"sub"`
	Nome           *string     `json:"nome"`
	Sobrenome      *string     `json:"sobrenome"`
	Email          *string     `json:"email"`
	LingoEXP       *int        `json:"LingoEXP"`
	Level          *int        `json:"Level"`
	Gender         *string     `json:"gender"`
	DataNascimento *string     `json:"data_nascimento"`
	Tokens         *int        `json:"tokens"`
	Gemas          *int        `json:"gemas"`
	Battery        *int        `json:"battery"`
	Plano          *string     `json:"plano"`
	Ranking        *int        `json:"ranking"`
	Listening      *int        `json:"listening"`
	Writing        *int        `json:"writing"`
	Reading        *int        `json:"reading"`
	Speaking       *int        `json:"speaking"`
	Difficulty     *string     `json:"difficulty"`
	Learning       *string     `json:"learning"`
	Items          interface{} `json:"items"`
	DailyMissions  interface{} `json:"dailyMissions"`
	Achievements   interface{} `json:"achievements"`
}

type UpdateUserDataResponse struct {
	Mensagem    string `json:"mensagem"`
	AccessToken string `json:"access_token"`
}

// UpdateUserData atualiza os dados do usuário nas tabelas normalizadas e gera um novo JWT
func UpdateUserData(req UpdateUserDataRequest) (*UpdateUserDataResponse, error) {
	// Pega o ID do usuário (prioriza 'id', depois 'sub')
	var userID int
	if req.ID != nil {
		userID = *req.ID
	} else if req.Sub != nil {
		userID = *req.Sub
	} else {
		return nil, errors.New("ID do usuário não fornecido")
	}

	// Busca o usuário completo no banco
	usuarioCompleto, err := repositories.GetUsuarioByID(userID)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	// Validação de battery (0 a 10)
	if req.Battery != nil {
		battery := *req.Battery
		if battery > 10 {
			battery = 10
		} else if battery < 0 {
			battery = 0
		}
		req.Battery = &battery
	}

	// Atualiza campos da tabela usuario
	if req.Nome != nil {
		usuarioCompleto.Usuario.Nome = *req.Nome
	}
	if req.Sobrenome != nil {
		usuarioCompleto.Usuario.Sobrenome = req.Sobrenome
	}
	if req.Email != nil {
		usuarioCompleto.Usuario.Email = *req.Email
	}
	if req.Gender != nil {
		usuarioCompleto.Usuario.Gender = req.Gender
	}
	if req.DataNascimento != nil {
		usuarioCompleto.Usuario.DataNascimento = req.DataNascimento
	}

	// Atualiza campos da tabela usuario_economia
	if req.Tokens != nil {
		usuarioCompleto.Economia.Tokens = *req.Tokens
	}
	if req.Gemas != nil {
		usuarioCompleto.Economia.Gemas = *req.Gemas
	}
	if req.Battery != nil {
		usuarioCompleto.Economia.Battery = *req.Battery
	}
	if req.Plano != nil {
		usuarioCompleto.Economia.Plano = *req.Plano
	}

	// Atualiza campos da tabela usuario_progresso
	if req.LingoEXP != nil {
		usuarioCompleto.Progresso.LingoEXP = *req.LingoEXP
	}
	if req.Level != nil {
		usuarioCompleto.Progresso.Level = *req.Level
	}
	if req.Listening != nil {
		usuarioCompleto.Progresso.Listening = *req.Listening
	}
	if req.Writing != nil {
		usuarioCompleto.Progresso.Writing = *req.Writing
	}
	if req.Reading != nil {
		usuarioCompleto.Progresso.Reading = *req.Reading
	}
	if req.Speaking != nil {
		usuarioCompleto.Progresso.Speaking = *req.Speaking
	}
	if req.Ranking != nil {
		usuarioCompleto.Progresso.Ranking = *req.Ranking
	}
	if req.Difficulty != nil {
		usuarioCompleto.Progresso.Difficulty = *req.Difficulty
	}
	if req.Learning != nil {
		usuarioCompleto.Progresso.Learning = *req.Learning
	}

	// Atualiza campos da tabela usuario_conteudo
	if req.Items != nil {
		usuarioCompleto.Conteudo.Items = req.Items
	}
	if req.DailyMissions != nil {
		usuarioCompleto.Conteudo.DailyMissions = req.DailyMissions
	}
	if req.Achievements != nil {
		usuarioCompleto.Conteudo.Achievements = req.Achievements
	}

	// Atualiza no banco de dados (todas as tabelas necessárias)
	err = repositories.UpdateUsuarioCompleto(usuarioCompleto)
	if err != nil {
		return nil, errors.New("erro ao atualizar usuário")
	}

	// Monta os dados para o novo JWT (apenas o ID)
	userData := map[string]interface{}{
		"id": usuarioCompleto.Usuario.ID,
	}

	// Gera novo access token
	accessToken, err := utils.GenerateAccessToken(usuarioCompleto.Usuario.ID, userData)
	if err != nil {
		return nil, errors.New("erro ao gerar novo token")
	}

	return &UpdateUserDataResponse{
		Mensagem:    "Novo JWT gerado e usuário atualizado com sucesso!",
		AccessToken: accessToken,
	}, nil
}
