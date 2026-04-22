package auth

import (
	"testing"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	username := "samuel"
	userID := "12345"
	token, err := GenerateJWT(userID, username)
	if err != nil {
		t.Fatalf("error generating JWT: %v", err)
	}
	if token == "" {
		t.Fatal("generated token is empty")
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("error validating JWT token: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.Username != username {
		t.Errorf("expected userID %s, got %s", username, claims.Username)
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	_, err := ValidateToken("")
	if err == nil {
		t.Error("expected error for empty token, got nil")
	}
}

func TestValidateToken_TamperedToken(t *testing.T) {
	username := "samuel"
	userID := "12345"
	token, err := GenerateJWT(userID, username)
	if err != nil {
		t.Fatalf("error generating JWT: %v", err)
	}

	if len(token) < 1 {
		t.Fatal("generated token is too short to tamper")
	}
	tampered := token[:len(token)-1] + "x"

	_, err = ValidateToken(tampered)
	if err == nil {
		t.Error("expected error for tampered token, got nil")
	}
}
