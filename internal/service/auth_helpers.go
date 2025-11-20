package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"kerjakuy/internal/dto"
)

func buildAuthResponse(user *dto.UserDTO, tokens *dto.AuthTokens) *dto.AuthResponse {
	return &dto.AuthResponse{
		User:   *user,
		Tokens: *tokens,
	}
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
