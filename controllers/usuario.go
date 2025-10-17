package controllers

import (
	"lingobotAPI-GO/repositories"
	"lingobotAPI-GO/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsuarios(c *gin.Context) {
	usuarios, err := repositories.GetAllUsuarios()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar usuários"})
		return
	}
	c.JSON(http.StatusOK, usuarios)
}

func CriarUsuario(c *gin.Context) {
	var req services.CriarUsuarioRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	err := services.CriarUsuario(req)
	if err != nil {
		// Determina o status code baseado no tipo de erro
		statusCode := http.StatusBadRequest

		if err.Error() == "e-mail já cadastrado" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensagem": "Usuário criado com sucesso!"})
}

func Login(c *gin.Context) {
	var req services.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	response, err := services.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func UpdateUserData(c *gin.Context) {
	var req services.UpdateUserDataRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	response, err := services.UpdateUserData(req)
	if err != nil {
		statusCode := http.StatusBadRequest

		if err.Error() == "usuário não encontrado" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUsuarioProfile retorna dados básicos do perfil (nome, email, etc)
func GetUsuarioProfile(c *gin.Context) {
	idParam := c.Param("id")
	usuarioID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	profile, err := repositories.GetUsuarioProfile(usuarioID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Usuário não encontrado"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetUsuarioContent retorna economia, progresso e conteúdo
func GetUsuarioContent(c *gin.Context) {
	idParam := c.Param("id")
	usuarioID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	content, err := repositories.GetUsuarioContent(usuarioID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Dados não encontrados"})
		return
	}

	c.JSON(http.StatusOK, content)
}

// GetUsuarioSocial retorna dados sociais (referal_code, invited_by)
func GetUsuarioSocial(c *gin.Context) {
	idParam := c.Param("id")
	usuarioID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	social, err := repositories.GetUsuarioSocial(usuarioID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Dados sociais não encontrados"})
		return
	}

	c.JSON(http.StatusOK, social)
}

// GetUsuarioSecurity retorna dados de segurança (OTP) - ADMIN ONLY
func GetUsuarioSecurity(c *gin.Context) {
	idParam := c.Param("id")
	usuarioID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return
	}

	// TODO: Adicionar validação se o usuário é admin
	// Por enquanto, qualquer um com JWT válido pode acessar

	security, err := repositories.GetUsuarioSecurity(usuarioID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Dados de segurança não encontrados"})
		return
	}

	c.JSON(http.StatusOK, security)
}
