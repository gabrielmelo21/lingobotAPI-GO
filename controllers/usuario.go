package controllers

import (
	"lingobotAPI-GO/repositories"
	"lingobotAPI-GO/services"
	"net/http"

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

	// Pega o user_id do token JWT
	userIDFromToken, _ := c.Get("user_id")

	// Validação extra: verifica se está tentando atualizar outro usuário
	if req.ID != nil && *req.ID != userIDFromToken.(int) {
		c.JSON(http.StatusForbidden, gin.H{"erro": "Você não pode atualizar dados de outro usuário"})
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
