package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"kerjakuy/internal/dto"
	"kerjakuy/internal/models"
	"kerjakuy/internal/repository"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

var ErrOAuthProviderNotConfigured = errors.New("oauth provider belum dikonfigurasi")

type AuthMetadata struct {
	UserAgent string
	IP        string
}

type AuthConfig struct {
	Secret          string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type AuthClaims struct {
	UserID    uuid.UUID
	Email     string
	ExpiresAt time.Time
}

type AuthService interface {
	Register(ctx context.Context, req dto.CreateUserRequest, meta AuthMetadata) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest, meta AuthMetadata) (*dto.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string, meta AuthMetadata) (*dto.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	ValidateAccessToken(token string) (*AuthClaims, error)
	BeginOAuth(ctx context.Context, provider, redirectURI string) (*dto.OAuthRedirectResponse, error)
	HandleOAuthCallback(ctx context.Context, provider, code, state string, meta AuthMetadata) (*dto.AuthResponse, error)
}

type authService struct {
	userSvc     UserService
	sessionRepo repository.UserSessionRepository
	tokens      tokenManager
}

func NewAuthService(userSvc UserService, sessionRepo repository.UserSessionRepository, cfg AuthConfig) AuthService {
	tokenMgr := &simpleTokenManager{
		secret:          []byte(cfg.Secret),
		issuer:          cfg.Issuer,
		accessTTL:       cfg.AccessTokenTTL,
		refreshTTL:      cfg.RefreshTokenTTL,
	}
	return &authService{
		userSvc:     userSvc,
		sessionRepo: sessionRepo,
		tokens:      tokenMgr,
	}
}

func (s *authService) Register(ctx context.Context, req dto.CreateUserRequest, meta AuthMetadata) (*dto.AuthResponse, error) {
	userDTO, err := s.userSvc.CreateWithPassword(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, err := s.issueTokens(ctx, userDTO.ID, userDTO.Email, meta)
	if err != nil {
		return nil, err
	}
	return buildAuthResponse(userDTO, resp), nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest, meta AuthMetadata) (*dto.AuthResponse, error) {
	user, err := s.userSvc.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	userDTO := dto.UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	resp, err := s.issueTokens(ctx, user.ID, user.Email, meta)
	if err != nil {
		return nil, err
	}
	return buildAuthResponse(&userDTO, resp), nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string, meta AuthMetadata) (*dto.AuthResponse, error) {
	claims, err := s.tokens.ValidateToken(refreshToken, tokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	session, err := s.sessionRepo.FindByTokenHash(ctx, hashToken(refreshToken))
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	if err := s.sessionRepo.DeleteByID(ctx, session.ID); err != nil {
		return nil, err
	}

	userDTO, err := s.userSvc.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	resp, err := s.issueTokens(ctx, claims.UserID, claims.Email, meta)
	if err != nil {
		return nil, err
	}

	return buildAuthResponse(userDTO, resp), nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	session, err := s.sessionRepo.FindByTokenHash(ctx, hashToken(refreshToken))
	if err != nil {
		return err
	}
	return s.sessionRepo.DeleteByID(ctx, session.ID)
}

func (s *authService) ValidateAccessToken(token string) (*AuthClaims, error) {
	return s.tokens.ValidateToken(token, tokenTypeAccess)
}

func (s *authService) BeginOAuth(ctx context.Context, provider, redirectURI string) (*dto.OAuthRedirectResponse, error) {
	state, err := generateState()
	if err != nil {
		return nil, err
	}
	if redirectURI == "" {
		redirectURI = "https://app.kerjakuy.local/auth/callback"
	}
	authURL := fmt.Sprintf("https://auth.%s.example.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s",
		provider,
		"CHANGE_ME",
		base64.URLEncoding.EncodeToString([]byte(redirectURI)),
		state,
	)
	return &dto.OAuthRedirectResponse{
		Provider:        provider,
		AuthorizationURL: authURL,
		State:           state,
	}, nil
}

func (s *authService) HandleOAuthCallback(ctx context.Context, provider, code, state string, meta AuthMetadata) (*dto.AuthResponse, error) {
	return nil, fmt.Errorf("%w: %s", ErrOAuthProviderNotConfigured, provider)
}

func (s *authService) issueTokens(ctx context.Context, userID uuid.UUID, email string, meta AuthMetadata) (*dto.AuthTokens, error) {
	claims := AuthClaims{
		UserID: userID,
		Email:  email,
	}
	accessToken, err := s.tokens.GenerateAccessToken(claims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.tokens.GenerateRefreshToken(claims)
	if err != nil {
		return nil, err
	}

	if err := s.createSession(ctx, userID, refreshToken, meta); err != nil {
		return nil, err
	}

	return &dto.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokens.AccessTTL().Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (s *authService) createSession(ctx context.Context, userID uuid.UUID, refreshToken string, meta AuthMetadata) error {
	userAgent := meta.UserAgent
	ip := meta.IP
	session := &models.UserSession{
		UserID:    userID,
		TokenHash: hashToken(refreshToken),
		ExpiresAt: time.Now().Add(s.tokens.RefreshTTL()),
	}
	if userAgent != "" {
		session.UserAgent = &userAgent
	}
	if ip != "" {
		session.IPAddress = &ip
	}
	return s.sessionRepo.Create(ctx, session)
}
