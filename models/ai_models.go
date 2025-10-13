package models

// AIRequest representa a requisição para os endpoints de IA
type AIRequest struct {
	Text    string `json:"text" binding:"required"`
	Mistral bool   `json:"mistral,omitempty"`
	Cohere  bool   `json:"cohere,omitempty"`
	Groq    bool   `json:"groq,omitempty"`
}

// AIResponse representa a resposta dos serviços de IA
type AIResponse struct {
	Response string  `json:"response,omitempty"`
	Error    string  `json:"error,omitempty"`
	Time     float64 `json:"time_seconds,omitempty"`
}

// BenchmarkResponse representa o resultado do benchmark
type BenchmarkResponse map[string]AIResponse
