package auth

import (
	"testing"
)

func TestGenerateAndVerifyPasswordHash(t *testing.T) {
	password := "mysecret"
	hash, err := GeneratePasswordHash(password)
	if err != nil {
		t.Fatalf("error generating hash: %v", err)
	}
	if !VerifyPasswordHash(password, hash) {
		t.Error("the password should be valid for the generated hash")
	}
	if VerifyPasswordHash("another", hash) {
		t.Error("an incorrect password shouldnt be valid")
	}
}
