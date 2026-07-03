package util

import "testing"

// TestEncryptDecryptAPIKeyRoundTrip 校验 AES 加解密往返一致
func TestEncryptDecryptAPIKeyRoundTrip(t *testing.T) {
	key := []byte("1234567890123456")
	plaintext := "sk-test-api-key-12345"

	encoded, err := EncryptAPIKey(plaintext, key)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	got, err := DecryptAPIKey(encoded, key)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if got != plaintext {
		t.Fatalf("decrypted = %q, want %q", got, plaintext)
	}
}

// TestDecryptAPIKeyWrongKey 校验错误密钥解密失败
func TestDecryptAPIKeyWrongKey(t *testing.T) {
	key := []byte("1234567890123456")
	wrongKey := []byte("abcdefghijklmnop")

	encoded, err := EncryptAPIKey("secret", key)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if _, err := DecryptAPIKey(encoded, wrongKey); err == nil {
		t.Fatal("expected decrypt error with wrong key")
	}
}
