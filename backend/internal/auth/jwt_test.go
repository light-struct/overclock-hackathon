package auth

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	secret := "test-secret"
	userID := int64(123)
	role := "student"
	duration := 24 * time.Hour

	token, err := GenerateToken(secret, userID, role, duration)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Error("Expected non-empty token")
	}
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret"
	userID := int64(123)
	role := "student"
	duration := 24 * time.Hour

	token, err := GenerateToken(secret, userID, role, duration)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}
}

func TestValidateTokenInvalid(t *testing.T) {
	secret := "test-secret"
	invalidToken := "invalid.token.here"

	_, err := ValidateToken(invalidToken, secret)
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestValidateTokenWrongSecret(t *testing.T) {
	secret := "test-secret"
	wrongSecret := "wrong-secret"
	userID := int64(123)
	role := "student"
	duration := 24 * time.Hour

	token, _ := GenerateToken(secret, userID, role, duration)

	_, err := ValidateToken(token, wrongSecret)
	if err == nil {
		t.Error("Expected error for wrong secret")
	}
}
