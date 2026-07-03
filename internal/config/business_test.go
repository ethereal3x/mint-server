package config

import "testing"

// TestValidateEmptyJWTSecret 校验 JWT secret 为空时返回错误
func TestValidateEmptyJWTSecret(t *testing.T) {
	cfg := BusinessConfig{
		JWT:       JWTConfig{Secret: "  "},
		SecretKey: SecretKeyCfg{Encryption: "1234567890123456"},
	}
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for empty jwt secret")
	}
}

// TestValidateInvalidAESKeyLength 校验 AES 密钥长度非法时返回错误
func TestValidateInvalidAESKeyLength(t *testing.T) {
	cfg := BusinessConfig{
		JWT:       JWTConfig{Secret: "test-secret"},
		SecretKey: SecretKeyCfg{Encryption: "short"},
	}
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for invalid aes key length")
	}
}

// TestValidateSuccess 校验合法配置通过校验
func TestValidateSuccess(t *testing.T) {
	cfg := BusinessConfig{
		JWT:       JWTConfig{Secret: "test-secret"},
		SecretKey: SecretKeyCfg{Encryption: "1234567890123456"},
	}
	if err := Validate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
