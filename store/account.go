package store

import (
	"context"
	"time"

	"github.com/twinj/uuid"
)

type AccountRepository interface {
	ResetFailedLoginCounter(userID uuid.UUID) error
	UpdateLoginAttempt(loggedUserID uuid.UUID, minutes float64, maxLoginFailures int) (int64, error)
	AddLoginToken(userID uuid.UUID, expirationTime int64, token, tokenType string) error
	SetPassword(userID uuid.UUID, password, token string) (*User, error)
	SetNewPassword(userID uuid.UUID, password string) error
	GetUserRoles(userID uuid.UUID) ([]*Role, error)
}

type AccountInMemRepository interface {
	SetSession(ctx context.Context, uid uuid.UUID, sid, v string, redisTokenTTL time.Duration) error
	GetSession(ctx context.Context, uid uuid.UUID, sid string) (string, error)
	DelSession(ctx context.Context, uid uuid.UUID, sid string) error
	DelSessionWithKey(ctx context.Context, key string) error
}

// GetRoles Get all available Roles.
func GetRoles() Roles {
	return Roles{
		User:  Role{Name: "User", Value: "USER"},
		Admin: Role{Name: "Admin", Value: "ADMIN"},
	}
}

// Roles object contains all roles.
type Roles struct {
	User  Role
	Admin Role
}

// Role is object with name and value.
type Role struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	ID    int64  `json:"id"`
}

// RolesDataResponse contains list of all roles.
type RolesDataResponse struct {
	Roles []Role `json:"roles"`
}
