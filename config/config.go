package config

import (
	"os"
	"time"

	"github.com/adinovcina/golang-setup/tools/logger"
	"github.com/joho/godotenv"
)

const (
	// Service default fallback values.
	stageDevelopment = "dev"
	apiPortDefault   = "5500"
	logLevelInfo     = "1"

	// Database default fallback values.
	migrationEnabledDefault = true

	// TTL default fallback values.
	mfaTemporaryTokenExpirationDefault = 5 * time.Minute
	mfaAccessTokenExpirationDefault    = 24 * time.Hour
	mfaRefreshTokenExpirationDefault   = 30 * 24 * time.Hour
	redisTokenExpirationDefault        = 24 * time.Hour

	maxLoginFailures       = 10
	banDurationDefaultTime = 5 * time.Minute
	timeoutDuration        = 30 * time.Second
)

// Load application configuration.
func Load() (config *Config, err error) {
	// load .env file
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		logger.Warn().Err(err).Msg(".env file not presented. Retrieving configuration from environment variables")
	}

	return loadFromEnv()
}
