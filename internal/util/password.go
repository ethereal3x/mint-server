package util

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const minPasswordLength = 8

// ErrPasswordTooShort 表示密码长度不足
var ErrPasswordTooShort = errors.New("password too short")

// HashPassword 生成 bcrypt 密码哈希
func HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword 校验明文密码和哈希是否匹配
func ComparePassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// ValidatePassword 校验密码基础强度
func ValidatePassword(password string) error {
	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}
	return nil
}
