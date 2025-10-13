package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// CallGemini chama a API do Google Gemini
func CallGemini(text string) (string, error) {
	apiKey := os.Getenv("GOOGLE_GEMINI_API_KEY1")
	if apiKey == "" {
		return "", errors.New("gemini API key not configured")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": text},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini API returned status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// Navega pela estrutura de resposta do Gemini
	if candidates, ok := result["candidates"].([]interface{}); ok && len(candidates) > 0 {
		if candidate, ok := candidates[0].(map[string]interface{}); ok {
			if content, ok := candidate["content"].(map[string]interface{}); ok {
				if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
					if part, ok := parts[0].(map[string]interface{}); ok {
						if textContent, ok := part["text"].(string); ok {
							return textContent, nil
						}
					}
				}
			}
		}
	}

	return "", errors.New("no text found in Gemini response")
}

// CallMistral chama a API do Mistral com retry
func CallMistral(text string) (string, error) {
	apiKey := os.Getenv("MISTRAL_KEY")
	if apiKey == "" {
		return "", errors.New("mistral API key not configured")
	}

	url := "https://api.mistral.ai/v1/chat/completions"
	maxRetries := 3

	payload := map[string]interface{}{
		"model": "mistral-tiny",
		"messages": []map[string]string{
			{"role": "user", "content": text},
		},
		"temperature": 0.7,
		"max_tokens":  2000,
	}

	jsonData, _ := json.Marshal(payload)

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			if attempt < maxRetries-1 {
				time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
				continue
			}
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 429 && attempt < maxRetries-1 {
			time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("mistral API returned status %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if message, ok := choice["message"].(map[string]interface{}); ok {
					if content, ok := message["content"].(string); ok {
						return content, nil
					}
				}
			}
		}

		return "", errors.New("no text found in Mistral response")
	}

	return "", errors.New("mistral request failed after retries")
}

// CallCohere chama a API do Cohere
func CallCohere(text string) (string, error) {
	apiKey := os.Getenv("COHERE_KEY")
	if apiKey == "" {
		return "", errors.New("cohere API key not configured")
	}

	url := "https://api.cohere.ai/v1/chat"

	payload := map[string]interface{}{
		"message":     text,
		"model":       "command-r",
		"temperature": 0.7,
		"max_tokens":  1000,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cohere API returned status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if textContent, ok := result["text"].(string); ok {
		return textContent, nil
	}

	return "", errors.New("no text found in Cohere response")
}

// CallGroq chama a API do Groq (via OpenAI SDK format)
func CallGroq(text string) (string, error) {
	apiKey := os.Getenv("GROQ_KEY")
	if apiKey == "" {
		return "", errors.New("groq API key not configured")
	}

	url := "https://api.groq.com/openai/v1/chat/completions"

	payload := map[string]interface{}{
		"model": "meta-llama/llama-4-scout-17b-16e-instruct",
		"messages": []map[string]string{
			{"role": "user", "content": text},
		},
		"temperature": 0.7,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq API returned status %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	return "", errors.New("no text found in Groq response")
}

// CallOpenRouter chama a API do OpenRouter com fallback de modelos
func CallOpenRouter(text string) (string, error) {
	apiKey := os.Getenv("OPENROUTER_KEY")
	if apiKey == "" {
		return "", errors.New("openRouter API key not configured")
	}

	url := "https://openrouter.ai/api/v1/chat/completions"

	modelsToTry := []string{
		"qwen/qwen3-235b-a22b-07-25:free",
		"meta-llama/llama-3.1-8b-instruct:free",
		"microsoft/phi-3-mini-128k-instruct:free",
		"google/gemma-2-9b-it:free",
	}

	for _, model := range modelsToTry {
		payload := map[string]interface{}{
			"model": model,
			"messages": []map[string]string{
				{"role": "user", "content": text},
			},
			"max_tokens":  1000,
			"temperature": 0.7,
		}

		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("HTTP-Referer", "https://lingobot-api.onrender.com")
		req.Header.Set("X-Title", "Go Gin OpenRouter App")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			var result map[string]interface{}
			json.Unmarshal(body, &result)

			if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
				if choice, ok := choices[0].(map[string]interface{}); ok {
					if message, ok := choice["message"].(map[string]interface{}); ok {
						if content, ok := message["content"].(string); ok {
							return content, nil
						}
					}
				}
			}
		}

		if resp.StatusCode == 503 {
			continue
		}
	}

	return "", errors.New("todos os modelos estão indisponíveis no momento")
}

// CallAIWithFallback tenta todas as IAs em ordem até obter sucesso
func CallAIWithFallback(text string, forceMistral, forceCohere, forceGroq bool) (string, error) {
	// Se forçar Mistral
	if forceMistral {
		return CallMistral(text)
	}

	// Se forçar Cohere
	if forceCohere {
		return CallCohere(text)
	}

	// Se forçar Groq
	if forceGroq {
		return CallGroq(text)
	}

	// Fallback: tenta todas em ordem
	// 1. Gemini
	if response, err := CallGemini(text); err == nil {
		return response, nil
	}

	// 2. Mistral
	if response, err := CallMistral(text); err == nil {
		return response, nil
	}

	// 3. Cohere
	if response, err := CallCohere(text); err == nil {
		return response, nil
	}

	// 4. Groq
	if response, err := CallGroq(text); err == nil {
		return response, nil
	}

	// 5. OpenRouter
	if response, err := CallOpenRouter(text); err == nil {
		return response, nil
	}

	return "", errors.New("all AI services failed")
}
