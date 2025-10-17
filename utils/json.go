package utils

import (
	"github.com/bytedance/sonic"
)

// Marshal converte struct para JSON usando Sonic
func Marshal(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

// Unmarshal converte JSON para struct usando Sonic
func Unmarshal(data []byte, v interface{}) error {
	return sonic.Unmarshal(data, v)
}

// MarshalString converte struct para JSON string usando Sonic
func MarshalString(v interface{}) (string, error) {
	return sonic.MarshalString(v)
}

// UnmarshalString converte JSON string para struct usando Sonic
func UnmarshalString(data string, v interface{}) error {
	return sonic.UnmarshalString(data, v)
}
