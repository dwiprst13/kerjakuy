package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"kerjakuy/internal/models"
	"kerjakuy/internal/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

var ErrOAuthProviderNotConfigured = errors.New("oauth provider belum dikonfigurasi")

type Metadata struct {
	UserAgent string
	IP        string
}

type Config struct {
	Secret          string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Claims struct {
	UserID    uuid.UUID
	Email     string
	ExpiresAt time.Time
}

type Service interface {
	Register(ctx context.Context, req user.CreateUserRequest, meta Metadata) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest, meta Metadata) (*AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string, meta Metadata) (*AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	ValidateAccessToken(token string) (*Claims, error)
	BeginOAuth(ctx context.Context, provider, redirectURI string) (*OAuthRedirectResponse, error)
	HandleOAuthCallback(ctx context.Context, provider, code, state string, meta Metadata) (*AuthResponse, error)
}

type userManager interface {
	CreateWithPassword(ctx context.Context, req user.CreateUserRequest) (*user.UserDTO, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user.UserDTO, error)
}

type authService struct {
	userSvc     userManager
	sessionRepo UserSessionRepository
	tokens      tokenManager
}

func NewService(userSvc userManager, sessionRepo UserSessionRepository, cfg Config) Service {
	tokenMgr := &jwtTokenManager{
		secret:     []byte(cfg.Secret),
		issuer:     cfg.Issuer,
		accessTTL:  cfg.AccessTokenTTL,
		refreshTTL: cfg.RefreshTokenTTL,
	}
	return &authService{
		userSvc:     userSvc,
		sessionRepo: sessionRepo,
		tokens:      tokenMgr,
	}
}

func (s *authService) Register(ctx context.Context, req user.CreateUserRequest, meta Metadata) (*AuthResponse, error) {
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

func (s *authService) Login(ctx context.Context, req LoginRequest, meta Metadata) (*AuthResponse, error) {
	account, err := s.userSvc.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	userDTO := user.UserDTO{
		ID:        account.ID,
		Name:      account.Name,
		Email:     account.Email,
		AvatarURL: account.AvatarURL,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}

	resp, err := s.issueTokens(ctx, account.ID, account.Email, meta)
	if err != nil {
		return nil, err
	}
	return buildAuthResponse(&userDTO, resp), nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string, meta Metadata) (*AuthResponse, error) {
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

func (s *authService) ValidateAccessToken(token string) (*Claims, error) {
	return s.tokens.ValidateToken(token, tokenTypeAccess)
}

func (s *authService) BeginOAuth(ctx context.Context, provider, redirectURI string) (*OAuthRedirectResponse, error) {
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
	return &OAuthRedirectResponse{
		Provider:         provider,
		AuthorizationURL: authURL,
		State:            state,
	}, nil
}

func (s *authService) HandleOAuthCallback(ctx context.Context, provider, code, state string, meta Metadata) (*AuthResponse, error) {
	return nil, fmt.Errorf("%w: %s", ErrOAuthProviderNotConfigured, provider)
}

func (s *authService) issueTokens(ctx context.Context, userID uuid.UUID, email string, meta Metadata) (*AuthTokens, error) {
	claims := Claims{
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

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokens.AccessTTL().Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (s *authService) createSession(ctx context.Context, userID uuid.UUID, refreshToken string, meta Metadata) error {
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
