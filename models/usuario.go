package models

import "time"

// Usuario - Tabela principal com dados básicos
type Usuario struct {
	ID             int       `json:"id" db:"id"`
	Nome           string    `json:"nome" db:"nome"`
	Sobrenome      *string   `json:"sobrenome" db:"sobrenome"`
	Email          string    `json:"email" db:"email"`
	Password       string    `json:"-" db:"password"` // nunca expor no JSON
	Gender         *string   `json:"gender" db:"gender"`
	DataNascimento *string   `json:"data_nascimento" db:"data_nascimento"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// UsuarioSeguranca - Dados de segurança
type UsuarioSeguranca struct {
	ID        int       `json:"id" db:"id"`
	UsuarioID int       `json:"usuario_id" db:"usuario_id"`
	OTPCode   *string   `json:"otp_code" db:"otp_code"`
	OTPAtivo  bool      `json:"otp_ativo" db:"otp_ativo"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UsuarioEconomia - Moedas, tokens e recursos
type UsuarioEconomia struct {
	ID        int       `json:"id" db:"id"`
	UsuarioID int       `json:"usuario_id" db:"usuario_id"`
	Tokens    int       `json:"tokens" db:"tokens"`
	Gemas     int       `json:"gemas" db:"gemas"`
	Battery   int       `json:"battery" db:"battery"`
	Plano     string    `json:"plano" db:"plano"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UsuarioProgresso - XP, level, skills
type UsuarioProgresso struct {
	ID         int       `json:"id" db:"id"`
	UsuarioID  int       `json:"usuario_id" db:"usuario_id"`
	LingoEXP   int       `json:"LingoEXP" db:"lingo_exp"`
	Level      int       `json:"Level" db:"level"`
	Listening  int       `json:"listening" db:"listening"`
	Writing    int       `json:"writing" db:"writing"`
	Reading    int       `json:"reading" db:"reading"`
	Speaking   int       `json:"speaking" db:"speaking"`
	Ranking    int       `json:"ranking" db:"ranking"`
	Difficulty string    `json:"difficulty" db:"difficulty"`
	Learning   string    `json:"learning" db:"learning"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// UsuarioSocial - Códigos de referência e convites
type UsuarioSocial struct {
	ID          int       `json:"id" db:"id"`
	UsuarioID   int       `json:"usuario_id" db:"usuario_id"`
	ReferalCode *string   `json:"referal_code" db:"referal_code"`
	InvitedBy   *string   `json:"invited_by" db:"invited_by"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UsuarioConteudo - Items, missions, achievements (JSONB)
type UsuarioConteudo struct {
	ID            int         `json:"id" db:"id"`
	UsuarioID     int         `json:"usuario_id" db:"usuario_id"`
	Items         interface{} `json:"items" db:"items"`
	DailyMissions interface{} `json:"dailyMissions" db:"daily_missions"`
	Achievements  interface{} `json:"achievements" db:"achievements"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}

// UsuarioCompleto - DTO com todos os dados
type UsuarioCompleto struct {
	Usuario   Usuario           `json:"usuario"`
	Seguranca *UsuarioSeguranca `json:"seguranca,omitempty"`
	Economia  UsuarioEconomia   `json:"economia"`
	Progresso UsuarioProgresso  `json:"progresso"`
	Social    *UsuarioSocial    `json:"social,omitempty"`
	Conteudo  *UsuarioConteudo  `json:"conteudo,omitempty"`
}

// UsuarioResponse - Resposta pública (sem dados sensíveis)
type UsuarioResponse struct {
	ID             int              `json:"id"`
	Nome           string           `json:"nome"`
	Sobrenome      *string          `json:"sobrenome"`
	Gender         *string          `json:"gender"`
	DataNascimento *string          `json:"data_nascimento"`
	CreatedAt      time.Time        `json:"created_at"`
	Economia       UsuarioEconomia  `json:"economia"`
	Progresso      UsuarioProgresso `json:"progresso"`
	Social         *UsuarioSocial   `json:"social,omitempty"`
	Conteudo       *UsuarioConteudo `json:"conteudo,omitempty"`
}

// UsuarioContentResponse - Resposta com economia, progresso e conteúdo
type UsuarioContentResponse struct {
	Economia  UsuarioEconomia  `json:"economia"`
	Progresso UsuarioProgresso `json:"progresso"`
	Conteudo  UsuarioConteudo  `json:"conteudo"`
}
