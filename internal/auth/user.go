package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"real-time-chat/internal/db"
	"real-time-chat/internal/models"

	"golang.org/x/crypto/argon2"
)

func GeneratePasswordHash(password string) (string, error) {
	salt := make([]byte, 16)

	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return fmt.Sprintf("%s$%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash)), nil
}

func VerifyPasswordHash(password, encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 2 {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	testHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return subtle.ConstantTimeCompare(hash, testHash) == 1
}

func AuthenticateUser(username, password string) (*models.User, error) {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	if !VerifyPasswordHash(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}
	return user, nil
}
