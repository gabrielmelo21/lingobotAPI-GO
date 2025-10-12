package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// getJWTSecret obtém a chave secreta do .env
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		panic("JWT_SECRET_KEY não está definida no .env")
	}
	return []byte(secret)
}

// Claims customizado para manter compatibilidade com Flask-JWT
type Claims struct {
	Fresh          bool        `json:"fresh"`
	JTI            string      `json:"jti"`
	Type           string      `json:"type"`
	Sub            int         `json:"sub"`
	CSRF           string      `json:"csrf"`
	ID             int         `json:"id"`
	Nome           string      `json:"nome"`
	Sobrenome      interface{} `json:"sobrenome"`
	Email          string      `json:"email"`
	LingoEXP       int         `json:"LingoEXP"`
	Level          int         `json:"Level"`
	Gender         interface{} `json:"gender"`
	DataNascimento interface{} `json:"data_nascimento"`
	Tokens         int         `json:"tokens"`
	Plano          string      `json:"plano"`
	CreatedAt      string      `json:"created_at"`
	ReferalCode    interface{} `json:"referal_code"`
	InvitedBy      interface{} `json:"invited_by"`
	Ranking        int         `json:"ranking"`
	Listening      int         `json:"listening"`
	Writing        int         `json:"writing"`
	Reading        int         `json:"reading"`
	Speaking       int         `json:"speaking"`
	Gemas          int         `json:"gemas"`
	Items          interface{} `json:"items"`
	DailyMissions  interface{} `json:"dailyMissions"`
	Achievements   interface{} `json:"achievements"`
	Difficulty     string      `json:"difficulty"`
	Battery        int         `json:"battery"`
	Learning       string      `json:"learning"`
	jwt.RegisteredClaims
}

// GenerateAccessToken gera um token de acesso JWT compatível com Flask-JWT
func GenerateAccessToken(userID int, userData map[string]interface{}) (string, error) {
	now := time.Now()
	exp := now.Add(7 * 24 * time.Hour)

	claims := Claims{
		Fresh:          false,
		JTI:            uuid.New().String(),
		Type:           "access",
		Sub:            userID,
		CSRF:           uuid.New().String(),
		ID:             userID,
		Nome:           getString(userData, "nome"),
		Sobrenome:      userData["sobrenome"],
		Email:          getString(userData, "email"),
		LingoEXP:       getInt(userData, "lingo_exp"),
		Level:          getInt(userData, "level"),
		Gender:         userData["gender"],
		DataNascimento: userData["data_nascimento"],
		Tokens:         getInt(userData, "tokens"),
		Plano:          getString(userData, "plano"),
		CreatedAt:      getString(userData, "created_at"),
		ReferalCode:    userData["referal_code"],
		InvitedBy:      userData["invited_by"],
		Ranking:        getInt(userData, "ranking"),
		Listening:      getInt(userData, "listening"),
		Writing:        getInt(userData, "writing"),
		Reading:        getInt(userData, "reading"),
		Speaking:       getInt(userData, "speaking"),
		Gemas:          getInt(userData, "gemas"),
		Items:          userData["items"],
		DailyMissions:  userData["daily_missions"],
		Achievements:   userData["achievements"],
		Difficulty:     getString(userData, "difficulty"),
		Battery:        getInt(userData, "battery"),
		Learning:       getString(userData, "learning"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// GenerateRefreshToken gera um token de refresh com duração de 30 dias
func GenerateRefreshToken(userID int) (string, error) {
	now := time.Now()
	exp := now.Add(30 * 24 * time.Hour)

	claims := jwt.MapClaims{
		"fresh": false,
		"jti":   uuid.New().String(),
		"type":  "refresh",
		"sub":   userID,
		"csrf":  uuid.New().String(),
		"exp":   exp.Unix(),
		"iat":   now.Unix(),
		"nbf":   now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ValidateToken valida e decodifica um token JWT
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}

// Funções auxiliares para converter tipos
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if val, ok := m[key].(int); ok {
		return val
	}
	return 0
}
