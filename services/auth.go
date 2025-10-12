package services

import (
	"errors"
	"lingobotAPI-GO/repositories"
	"lingobotAPI-GO/utils"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Mensagem     string `json:"mensagem"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Login realiza o login do usuário e retorna tokens JWT
func Login(req LoginRequest) (*LoginResponse, error) {
	// Busca o usuário pelo email
	usuario, err := repositories.GetUsuarioByEmail(req.Email)
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Verifica a senha
	if !utils.VerifyPassword(req.Password, usuario.Password) {
		return nil, errors.New("credenciais inválidas")
	}

	// Monta os dados do usuário para o token
	userData := map[string]interface{}{
		"id":              usuario.ID,
		"nome":            usuario.Nome,
		"sobrenome":       usuario.Sobrenome,
		"email":           usuario.Email,
		"password":        usuario.Password,
		"otp_code":        usuario.OTPCode,
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

	// Gera os tokens
	accessToken, err := utils.GenerateAccessToken(usuario.ID, userData)
	if err != nil {
		return nil, errors.New("erro ao gerar token de acesso")
	}

	refreshToken, err := utils.GenerateRefreshToken(usuario.ID)
	if err != nil {
		return nil, errors.New("erro ao gerar token de refresh")
	}

	return &LoginResponse{
		Mensagem:     "Login realizado com sucesso!",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
