package util

import (
	"context"
	"testing"
	"time"
)

// TestTokenManagerIssueAndParse 校验 JWT 签发与解析往返一致
func TestTokenManagerIssueAndParse(t *testing.T) {
	manager := NewTokenManager(TokenManagerConfig{
		Secret: "test-jwt-secret",
		TTL:    time.Hour,
	})
	result, err := manager.IssueAccessToken(context.Background(), &TokenInput{
		UserID:     "user-1",
		Provider:   "password",
		Identifier: "alice",
	})
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	claims, err := manager.ParseAccessToken(context.Background(), result.AccessToken)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.Subject != "user-1" {
		t.Fatalf("subject = %q, want user-1", claims.Subject)
	}
	if claims.Provider != "password" || claims.Identifier != "alice" {
		t.Fatalf("unexpected claims: provider=%q identifier=%q", claims.Provider, claims.Identifier)
	}
}

// TestTokenManagerParseExpired 校验过期令牌解析失败
func TestTokenManagerParseExpired(t *testing.T) {
	manager := NewTokenManager(TokenManagerConfig{
		Secret: "test-jwt-secret",
		TTL:    time.Millisecond,
	})
	result, err := manager.IssueAccessToken(context.Background(), &TokenInput{UserID: "user-1"})
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	time.Sleep(2 * time.Millisecond)
	if _, err := manager.ParseAccessToken(context.Background(), result.AccessToken); err == nil {
		t.Fatal("expected expired token error")
	}
}
