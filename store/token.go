package store

import "github.com/twinj/uuid"

const (
	TokenNotFound = "token not found"
)

type TokenRepository interface {
	AddPasswordResetToken(userID uuid.UUID, token string, expiresAt int64) (*PasswordToken, error)
	GetPasswordTokenByToken(token string) (*PasswordToken, error)
	GetTokenByTokenAndType(token, tokenType string) (*LoginToken, error)
	DeleteTokenByID(id int64) error
}

// GetTokenTypes get available token types.
func GetTokenTypes() Tokens {
	return Tokens{
		MFA:          "MFA",
		RefreshToken: "REFRESH_TOKEN",
	}
}

// Tokens struct used to describe token types.
type Tokens struct {
	MFA          string
	Authorize    string
	RefreshToken string
}

// LoginToken represents token struct.
type LoginToken struct {
	Token     string    `json:"token"`
	TokenType string    `json:"tokenType"`
	ID        int64     `json:"id"`
	UserID    uuid.UUID `json:"userID"`
	Expired   bool      `json:"expired"`
}

// PasswordToken contains innfo about password token.
type PasswordToken struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt int64
}
