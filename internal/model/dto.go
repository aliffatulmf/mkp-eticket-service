package model

type CreateTerminalRequest struct {
	Code    string `json:"code" validate:"required,max=10"`
	Name    string `json:"name" validate:"required,max=100"`
	Address string `json:"address" validate:"required"`
}

type UpdateTerminalRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Address  string `json:"address" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type CreateAdminRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateAdminRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminLoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	ExpiresAt    int64  `json:"expires_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}
