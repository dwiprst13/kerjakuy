package auth

import "kerjakuy/internal/user"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User   user.UserDTO `json:"user"`
	Tokens AuthTokens   `json:"tokens"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type OAuthRedirectResponse struct {
	Provider         string `json:"provider"`
	AuthorizationURL string `json:"authorization_url"`
	State            string `json:"state"`
}
