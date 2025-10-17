package services

import (
	"errors"
	"fmt"
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
	fmt.Printf("🔍 Tentando login para: %s\n", req.Email)

	// Busca o usuário completo pelo email
	usuarioCompleto, err := repositories.GetUsuarioByEmail(req.Email)
	if err != nil {
		fmt.Printf("❌ Erro ao buscar usuário: %v\n", err)
		return nil, errors.New("credenciais inválidas")
	}

	if usuarioCompleto == nil {
		fmt.Printf("❌ Usuário não encontrado\n")
		return nil, errors.New("credenciais inválidas")
	}

	fmt.Printf("✅ Usuário encontrado: %s (ID: %d)\n", usuarioCompleto.Usuario.Email, usuarioCompleto.Usuario.ID)

	// Verifica a senha
	senhaCorreta := utils.VerifyPassword(req.Password, usuarioCompleto.Usuario.Password)
	fmt.Printf("🔑 Senha correta: %v\n", senhaCorreta)

	if !senhaCorreta {
		fmt.Printf("❌ Senha incorreta\n")
		return nil, errors.New("credenciais inválidas")
	}

	fmt.Printf("✅ Senha validada com sucesso\n")

	// JWT minimalista - apenas o ID do usuário
	// O frontend deve buscar os dados completos após o login se necessário
	userData := map[string]interface{}{
		"id": usuarioCompleto.Usuario.ID,
	}

	// Gera os tokens
	accessToken, err := utils.GenerateAccessToken(usuarioCompleto.Usuario.ID, userData)
	if err != nil {
		fmt.Printf("❌ Erro ao gerar access token: %v\n", err)
		return nil, errors.New("erro ao gerar token de acesso")
	}

	refreshToken, err := utils.GenerateRefreshToken(usuarioCompleto.Usuario.ID)
	if err != nil {
		fmt.Printf("❌ Erro ao gerar refresh token: %v\n", err)
		return nil, errors.New("erro ao gerar token de refresh")
	}

	fmt.Printf("✅ Tokens gerados com sucesso\n")

	return &LoginResponse{
		Mensagem:     "Login realizado com sucesso!",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
