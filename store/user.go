package store

import (
	"time"

	"github.com/twinj/uuid"
)

const (
	UserNotFound   = "user not found"
	UserDuplicated = "duplicate user"
)

type UserRepository interface {
	GetUserByID(id uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByToken(token, tokenType string) (*User, error)
	ActivateUser(userID uuid.UUID) error
	GetUsers(filter *UserFilter) ([]*User, error)
	UpdateUser(user *User) (*User, error)
}

// User model.
type User struct {
	CreatedAt         time.Time  `json:"createdAt,omitempty"`
	LoginBlockedUntil *time.Time `json:"loginBlockedUntil,omitempty"`
	Name              string     `json:"name,omitempty"`
	Email             string     `json:"email,omitempty"`
	Phone             string     `json:"phone,omitempty"`
	Language          string     `json:"language,omitempty"`
	Password          string     `json:"-"`
	Role              string     `json:"role,omitempty"`
	RoleID            int64      `json:"roleID,omitempty"`
	FailedLoginCount  int        `json:"failedLoginCount,omitempty"`
	ID                uuid.UUID  `json:"id,omitempty"`
	Active            bool       `json:"active,omitempty"`
	EmailVerified     bool       `json:"emailVerified,omitempty"`
	TermsAccepted     bool       `json:"termsAccepted,omitempty"`
	Expired           bool       `json:"expired,omitempty"`
}

func (u *User) VerifyIfUserIsSuspended(maxLoginFailures int) bool {
	return u.LoginBlockedUntil != nil && u.LoginBlockedUntil.After(time.Now()) &&
		u.FailedLoginCount >= maxLoginFailures
}

type UserFilter struct {
	Active *bool
	Search *string
}

// UsersKeyToColumnMap - a map of sorting keys matching their DB values.
func UsersKeyToColumnMap() map[string]string {
	usersToColumnMap := map[string]string{
		"id": "u.id",
	}

	return usersToColumnMap
}
