package utils

import (
	"fmt"

	"github.com/twinj/uuid"
)

func GenerateUniqueID() string {
	return uuid.NewV4().String()
}

// FormatSessionKey - method generates session key for Redis.
func FormatSessionKey(userID uuid.UUID, sessionID string) string {
	return fmt.Sprintf("session:%v:%s", userID, sessionID)
}
