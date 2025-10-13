package models

// TTSRequest representa a requisição para text-to-speech
type TTSRequest struct {
	Text    string `json:"text" binding:"required"`
	Voice   int    `json:"voice"`
	Premium bool   `json:"premium"`
}

// TranscribeResponse representa a resposta da transcrição
type TranscribeResponse struct {
	Text  string `json:"text,omitempty"`
	Error string `json:"error,omitempty"`
}
