package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// DefaultCost определяет сложность хэширования (по умолчанию 10)
const DefaultCost = bcrypt.DefaultCost

// HashPassword создает безопасный хэш пароля с использованием bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword проверяет совпадает ли пароль с хэшем
func ComparePassword(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
