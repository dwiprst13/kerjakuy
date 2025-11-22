package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type tokenManager interface {
	GenerateAccessToken(claims Claims) (string, error)
	GenerateRefreshToken(claims Claims) (string, error)
	ValidateToken(token string, expectedType string) (*Claims, error)
	AccessTTL() time.Duration
	RefreshTTL() time.Duration
}

type simpleTokenManager struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type tokenPayload struct {
	UserID    string `json:"uid"`
	Email     string `json:"email"`
	Issuer    string `json:"iss"`
	TokenType string `json:"typ"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func (m *simpleTokenManager) GenerateAccessToken(claims Claims) (string, error) {
	return m.generateToken(claims, tokenTypeAccess, m.accessTTL)
}

func (m *simpleTokenManager) GenerateRefreshToken(claims Claims) (string, error) {
	return m.generateToken(claims, tokenTypeRefresh, m.refreshTTL)
}

func (m *simpleTokenManager) AccessTTL() time.Duration {
	return m.accessTTL
}

func (m *simpleTokenManager) RefreshTTL() time.Duration {
	return m.refreshTTL
}

func (m *simpleTokenManager) generateToken(claims Claims, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	payload := tokenPayload{
		UserID:    claims.UserID.String(),
		Email:     claims.Email,
		Issuer:    m.issuer,
		TokenType: tokenType,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(ttl).Unix(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	signature := m.sign(payloadBytes)
	return encodeSegment(payloadBytes) + "." + encodeSegment(signature), nil
}

func (m *simpleTokenManager) ValidateToken(token string, expectedType string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("token invalid")
	}

	payloadBytes, err := decodeSegment(parts[0])
	if err != nil {
		return nil, errors.New("token invalid")
	}

	signature, err := decodeSegment(parts[1])
	if err != nil {
		return nil, errors.New("token invalid")
	}

	if !hmac.Equal(signature, m.sign(payloadBytes)) {
		return nil, errors.New("token signature invalid")
	}

	var payload tokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, errors.New("token payload invalid")
	}

	if payload.TokenType != expectedType {
		return nil, errors.New("token type mismatch")
	}

	if payload.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return nil, errors.New("token payload invalid")
	}

	return &Claims{
		UserID:    userID,
		Email:     payload.Email,
		ExpiresAt: time.Unix(payload.ExpiresAt, 0),
	}, nil
}

func (m *simpleTokenManager) sign(payload []byte) []byte {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write(payload)
	return mac.Sum(nil)
}

func encodeSegment(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func decodeSegment(segment string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(segment)
}
