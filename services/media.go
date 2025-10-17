package services

import (
	"LingobotAPI-GO/utils"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var VoiceIDs = []string{
	"TxGEqnHWrfWFTfGW9XjX", // 0 - Josh
	"pNInz6obpgDQGcFmaJgB", // 1 - Adam
	"onwK4e9ZLuTAKqWW03F9", // 2 - James
	"yoZ06aMxZJJ28mfd3POQ", // 3 - Sam
	"VR6AewLTigWG4xSOukaG", // 4 - Arnold
	"EXAVITQu4vr4xnSDxMaL", // 5 - Bella (feminina padrão)
}

// GenerateTTSGoogle gera áudio usando edge-tts (Google TTS)
func GenerateTTSGoogle(text string) ([]byte, error) {
	// Cria arquivo temporário para o áudio
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("tts_%d.mp3", time.Now().UnixNano()))
	defer os.Remove(tempFile)

	// Executa edge-tts via comando
	cmd := exec.Command("edge-tts", "--text", text, "--voice", "en-US-ChristopherNeural", "--write-media", tempFile)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("erro ao executar edge-tts: %v", err)
	}

	// Lê o arquivo gerado
	audioData, err := os.ReadFile(tempFile)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de áudio: %v", err)
	}

	return audioData, nil
}

// GenerateTTSElevenLabs gera áudio usando ElevenLabs
func GenerateTTSElevenLabs(text, voiceID string) ([]byte, error) {
	apiKey := os.Getenv("ELEVENLABS_KEY1")
	if apiKey == "" {
		return nil, errors.New("ElevenLabs API key not configured")
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID)

	payload := map[string]interface{}{
		"text":     text,
		"model_id": "eleven_multilingual_v2",
		"voice_settings": map[string]interface{}{
			"stability":         0.5,
			"similarity_boost":  0.75,
			"style":             0.0,
			"use_speaker_boost": true,
		},
	}

	jsonData, _ := utils.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("xi-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "audio/mpeg")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ElevenLabs API returned status %d", resp.StatusCode)
	}

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return audioData, nil
}

// GenerateTTS gera áudio com fallback (ElevenLabs -> Google TTS)
func GenerateTTS(text string, voiceIndex int, premium bool) ([]byte, error) {
	// Valida índice de voz
	if voiceIndex < 0 || voiceIndex >= len(VoiceIDs) {
		voiceIndex = len(VoiceIDs) - 1 // Padrão: última voz
	}

	voiceID := VoiceIDs[voiceIndex]

	// Se premium, tenta ElevenLabs primeiro
	if premium {
		audioData, err := GenerateTTSElevenLabs(text, voiceID)
		if err == nil {
			return audioData, nil
		}
		fmt.Printf("⚠️ Falha com ElevenLabs: %v. Usando Google TTS como fallback...\n", err)
	}

	// Fallback: Google TTS
	return GenerateTTSGoogle(text)
}

// TranscribeAudio transcreve áudio usando AssemblyAI
func TranscribeAudio(filePath string) (string, error) {
	apiKey := os.Getenv("ASSEMBLYAI_KEY")
	if apiKey == "" {
		return "", errors.New("AssemblyAI API key not configured")
	}

	// 1. Upload do arquivo
	uploadURL := "https://api.assemblyai.com/v2/upload"

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	req, _ := http.NewRequest("POST", uploadURL, file)
	req.Header.Set("authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var uploadResp map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	utils.Unmarshal(body, &uploadResp)

	audioURL, ok := uploadResp["upload_url"].(string)
	if !ok {
		return "", errors.New("failed to get upload URL")
	}

	// 2. Solicita transcrição
	transcriptURL := "https://api.assemblyai.com/v2/transcript"

	transcriptPayload := map[string]interface{}{
		"audio_url":    audioURL,
		"speech_model": "universal",
	}

	jsonData, _ := utils.Marshal(transcriptPayload)
	req, _ = http.NewRequest("POST", transcriptURL, bytes.NewBuffer(jsonData))
	req.Header.Set("authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}

	var transcriptResp map[string]interface{}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	utils.Unmarshal(body, &transcriptResp)

	transcriptID, ok := transcriptResp["id"].(string)
	if !ok {
		return "", errors.New("failed to get transcript ID")
	}

	// 3. Polling: aguarda conclusão da transcrição
	pollURL := fmt.Sprintf("https://api.assemblyai.com/v2/transcript/%s", transcriptID)

	for {
		time.Sleep(2 * time.Second)

		req, _ = http.NewRequest("GET", pollURL, nil)
		req.Header.Set("authorization", apiKey)

		resp, err = client.Do(req)
		if err != nil {
			return "", err
		}

		var pollResp map[string]interface{}
		body, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		utils.Unmarshal(body, &pollResp)

		status, _ := pollResp["status"].(string)

		if status == "completed" {
			text, _ := pollResp["text"].(string)
			return text, nil
		}

		if status == "error" {
			errorMsg, _ := pollResp["error"].(string)
			return "", fmt.Errorf("transcription error: %s", errorMsg)
		}

		// Continua polling se status for "queued" ou "processing"
	}
}