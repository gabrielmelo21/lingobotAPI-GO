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

// Claims customizado - JWT minimalista com apenas dados essenciais
type Claims struct {
	Fresh bool   `json:"fresh"`
	JTI   string `json:"jti"`
	Type  string `json:"type"`
	Sub   int    `json:"sub"` // User ID
	CSRF  string `json:"csrf"`
	jwt.RegisteredClaims
}

// GenerateAccessToken gera um token de acesso JWT minimalista (7 dias)
func GenerateAccessToken(userID int, userData map[string]interface{}) (string, error) {
	now := time.Now()
	exp := now.Add(7 * 24 * time.Hour)

	claims := Claims{
		Fresh: false,
		JTI:   uuid.New().String(),
		Type:  "access",
		Sub:   userID, // Apenas o ID do usuário
		CSRF:  uuid.New().String(),
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
