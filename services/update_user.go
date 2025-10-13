package services

import (
	"errors"
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
	Plano          *string     `json:"plano"`
	Ranking        *int        `json:"ranking"`
	Listening      *int        `json:"listening"`
	Writing        *int        `json:"writing"`
	Reading        *int        `json:"reading"`
	Speaking       *int        `json:"speaking"`
	Gemas          *int        `json:"gemas"`
	Items          interface{} `json:"items"`
	DailyMissions  interface{} `json:"dailyMissions"`
	Achievements   interface{} `json:"achievements"`
	Difficulty     *string     `json:"difficulty"`
	Battery        *int        `json:"battery"`
	Learning       *string     `json:"learning"`
}

type UpdateUserDataResponse struct {
	Mensagem    string `json:"mensagem"`
	AccessToken string `json:"access_token"`
}

// UpdateUserData atualiza os dados do usuário e gera um novo JWT
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

	// Busca o usuário no banco
	usuario, err := repositories.GetUsuarioByID(userID)
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

	// Atualiza apenas os campos fornecidos
	if req.Nome != nil {
		usuario.Nome = *req.Nome
	}
	if req.Sobrenome != nil {
		usuario.Sobrenome = req.Sobrenome
	}
	if req.Email != nil {
		usuario.Email = *req.Email
	}
	if req.LingoEXP != nil {
		usuario.LingoEXP = *req.LingoEXP
	}
	if req.Level != nil {
		usuario.Level = *req.Level
	}
	if req.Gender != nil {
		usuario.Gender = req.Gender
	}
	if req.DataNascimento != nil {
		usuario.DataNascimento = req.DataNascimento
	}
	if req.Tokens != nil {
		usuario.Tokens = *req.Tokens
	}
	if req.Plano != nil {
		usuario.Plano = *req.Plano
	}
	if req.Ranking != nil {
		usuario.Ranking = *req.Ranking
	}
	if req.Listening != nil {
		usuario.Listening = *req.Listening
	}
	if req.Writing != nil {
		usuario.Writing = *req.Writing
	}
	if req.Reading != nil {
		usuario.Reading = *req.Reading
	}
	if req.Speaking != nil {
		usuario.Speaking = *req.Speaking
	}
	if req.Gemas != nil {
		usuario.Gemas = *req.Gemas
	}
	if req.Items != nil {
		usuario.Items = req.Items
	}
	if req.DailyMissions != nil {
		usuario.DailyMissions = req.DailyMissions
	}
	if req.Achievements != nil {
		usuario.Achievements = req.Achievements
	}
	if req.Difficulty != nil {
		usuario.Difficulty = *req.Difficulty
	}
	if req.Battery != nil {
		usuario.Battery = *req.Battery
	}
	if req.Learning != nil {
		usuario.Learning = *req.Learning
	}

	// Atualiza no banco de dados
	err = repositories.UpdateUsuario(usuario)
	if err != nil {
		return nil, errors.New("erro ao atualizar usuário")
	}

	// Monta os dados para o novo JWT
	userData := map[string]interface{}{
		"id":              usuario.ID,
		"nome":            usuario.Nome,
		"sobrenome":       usuario.Sobrenome,
		"email":           usuario.Email,
		"lingo_exp":       usuario.LingoEXP,
		"level":           usuario.Level,
		"gender":          usuario.Gender,
		"data_nascimento": usuario.DataNascimento,
		"tokens":          usuario.Tokens,
		"plano":           usuario.Plano,
		"created_at":      usuario.CreatedAt,
		"referal_code":    usuario.ReferalCode,
		"invited_by":      usuario.InvitedBy,
		"ranking":         usuario.Ranking,
		"listening":       usuario.Listening,
		"writing":         usuario.Writing,
		"reading":         usuario.Reading,
		"speaking":        usuario.Speaking,
		"gemas":           usuario.Gemas,
		"items":           usuario.Items,
		"daily_missions":  usuario.DailyMissions,
		"achievements":    usuario.Achievements,
		"difficulty":      usuario.Difficulty,
		"battery":         usuario.Battery,
		"learning":        usuario.Learning,
	}

	// Gera novo access token
	accessToken, err := utils.GenerateAccessToken(usuario.ID, userData)
	if err != nil {
		return nil, errors.New("erro ao gerar novo token")
	}

	return &UpdateUserDataResponse{
		Mensagem:    "Novo JWT gerado e usuário atualizado com sucesso!",
		AccessToken: accessToken,
	}, nil
}
