package controllers

import (
	"fmt"
	"lingobotAPI-GO/models"
	"lingobotAPI-GO/services"
	"lingobotAPI-GO/utils"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// TTS endpoint para text-to-speech
func TTS(c *gin.Context) {
	var req models.TTSRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SonicJSON(c, http.StatusBadRequest, gin.H{"error": "Texto é obrigatório"})
		return
	}

	text := req.Text
	voiceIndex := req.Voice
	premium := req.Premium

	// Valida índice de voz
	if voiceIndex < 0 || voiceIndex >= len(services.VoiceIDs) {
		utils.SonicJSON(c, http.StatusBadRequest, gin.H{"error": "Índice de voz inválido"})
		return
	}

	fmt.Printf("🔊 Gerando TTS para: %.60s... (voz %d) | Premium: %t\n", text, voiceIndex, premium)

	// Gera o áudio
	audioData, err := services.GenerateTTS(text, voiceIndex, premium)
	if err != nil {
		utils.SonicJSON(c, http.StatusInternalServerError, gin.H{"error": "Erro ao gerar áudio"})
		return
	}

	// Retorna o áudio como MP3
	c.Data(http.StatusOK, "audio/mp3", audioData)
}

// TranscribeAudio endpoint para transcrição de áudio
func TranscribeAudio(c *gin.Context) {
	// Verifica se foi enviado um arquivo
	file, err := c.FormFile("file")
	if err != nil {
		utils.SonicJSON(c, http.StatusBadRequest, gin.H{"error": "Nenhum arquivo enviado"})
		return
	}

	// Salva temporariamente o arquivo
	tempDir := os.TempDir()
	tempPath := filepath.Join(tempDir, fmt.Sprintf("audio_%d_%s", time.Now().UnixNano(), file.Filename))

	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		utils.SonicJSON(c, http.StatusInternalServerError, gin.H{"error": "Erro ao salvar arquivo"})
		return
	}
	defer os.Remove(tempPath) // Remove após processar

	fmt.Printf("📝 Transcrevendo áudio: %s\n", file.Filename)

	// Transcreve o áudio
	text, err := services.TranscribeAudio(tempPath)
	if err != nil {
		utils.SonicJSON(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retorna o texto transcrito
	utils.SonicJSON(c, http.StatusOK, models.TranscribeResponse{
		Text: text,
	})
}