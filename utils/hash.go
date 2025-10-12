package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashSenha cria um hash bcrypt da senha
func HashPassword(senha string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerificarSenha compara uma senha com seu hash
func VerifyPassword(senha, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha))
	return err == nil
}
