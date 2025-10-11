package models

type Usuario struct {
	ID             int     `json:"id"`
	Nome           string  `json:"nome"`
	Sobrenome      *string `json:"sobrenome"` // ponteiro para aceitar NULL
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	OTPCode        *string `json:"otp_code"`
	LingoEXP       int     `json:"lingo_exp"`
	Level          int     `json:"level"`
	Gender         *string `json:"gender"`
	DataNascimento *string `json:"data_nascimento"`
	Tokens         int     `json:"tokens"`
	Plano          string  `json:"plano"`
	CreatedAt      string  `json:"created_at"`
	ReferalCode    *string `json:"referal_code"`
	InvitedBy      *string `json:"invited_by"`
	Ranking        int     `json:"ranking"`
	Listening      int     `json:"listening"`
	Writing        int     `json:"writing"`
	Reading        int     `json:"reading"`
	Speaking       int     `json:"speaking"`
	Gemas          int     `json:"gemas"`
	Items          string  `json:"items"`
	DailyMissions  string  `json:"daily_missions"`
	Achievements   string  `json:"achievements"`
	Difficulty     string  `json:"difficulty"`
	Battery        int     `json:"battery"`
	Learning       string  `json:"learning"`
}
