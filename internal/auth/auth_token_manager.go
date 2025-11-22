package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type tokenManager interface {
	GenerateAccessToken(claims Claims) (string, error)
	GenerateRefreshToken(claims Claims) (string, error)
	ValidateToken(token string, expectedType string) (*Claims, error)
	AccessTTL() time.Duration
	RefreshTTL() time.Duration
}

type jwtTokenManager struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type jwtClaims struct {
	UserID    string `json:"uid"`
	Email     string `json:"email"`
	TokenType string `json:"typ"`
	jwt.RegisteredClaims
}

func (m *jwtTokenManager) GenerateAccessToken(claims Claims) (string, error) {
	return m.generateToken(claims, tokenTypeAccess, m.accessTTL)
}

func (m *jwtTokenManager) GenerateRefreshToken(claims Claims) (string, error) {
	return m.generateToken(claims, tokenTypeRefresh, m.refreshTTL)
}

func (m *jwtTokenManager) AccessTTL() time.Duration {
	return m.accessTTL
}

func (m *jwtTokenManager) RefreshTTL() time.Duration {
	return m.refreshTTL
}

func (m *jwtTokenManager) generateToken(claims Claims, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)

	jClaims := jwtClaims{
		UserID:    claims.UserID.String(),
		Email:     claims.Email,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   claims.UserID.String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jClaims)
	return token.SignedString(m.secret)
}

func (m *jwtTokenManager) ValidateToken(tokenString string, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		if claims.TokenType != expectedType {
			return nil, errors.New("invalid token type")
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			return nil, errors.New("invalid user id in token")
		}

		return &Claims{
			UserID:    userID,
			Email:     claims.Email,
			ExpiresAt: claims.ExpiresAt.Time,
		}, nil
	}

	return nil, errors.New("invalid token")
}
