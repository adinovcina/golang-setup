package api

import (
	"strings"
	"time"

	"github.com/adinovcina/golang-setup/tools/utils"
	jwt "github.com/golang-jwt/jwt"
	bg "github.com/kjk/betterguid"
	"github.com/twinj/uuid"
)

// Data contains basic user data after user is authorized.
type Data struct {
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	RequestID  string    `json:"requestID,omitempty"`
	SessionKey string    `json:"sessionKey"`
	UserID     uuid.UUID `json:"userID"`
	UserRoleID int64     `json:"userRoleID"`
	Active     bool      `json:"active"`
}

// Claim for JWT.
type Claim struct {
	jwt.StandardClaims
	SessionID string    `json:"sessionID"`
	UserID    uuid.UUID `json:"userID"`
}

// CreateToken Generate jwt new token for the user.
func (c *Claim) CreateToken(tokenTTL time.Duration, secretKey string) (string, error) {
	c.ExpiresAt = time.Now().Add(tokenTTL).Unix()
	c.StandardClaims.ExpiresAt = c.ExpiresAt

	// Create the token using your claims
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(secretKey))
}

// NewID will generate a new random ID.
func (c *Claim) NewID() {
	c.SessionID = bg.New()
}

// NewRefreshToken generates new UUID string token.
func NewRefreshToken() string {
	return NewDoubleUUIDCode()
}

func NewDoubleUUIDCode() string {
	return strings.ReplaceAll(utils.GenerateUniqueID()+utils.GenerateUniqueID(), "-", "")
}
