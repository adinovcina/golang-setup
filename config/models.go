package config

import (
	"time"
)

// Config stores application configuration.
type Config struct {
	Service  Service
	Database Database
	Redis    Redis
	Email    Email
	Timeouts Timeouts
	MFA      MFA
	Account  Account
}

// Service contains configuration for service.
type Service struct {
	Port        string
	Environment string
	LogLevel    string
}

// Account contains data related to login attempts.
type Account struct {
	MaxLoginFailures int
	BanDurationTime  time.Duration
}

// Timeouts contains configuration for read and write timeouts.
type Timeouts struct {
	ReadDuration  time.Duration
	WriteDuration time.Duration
}

// Database configuration.
type Database struct {
	Username         string
	Password         string
	Name             string
	Address          string
	Port             string
	MigrationFolder  string
	MigrationEnabled bool
}

// Redis stores configuration for connection to Redis database.
type Redis struct {
	Address   string
	Database  string
	Password  string
	SecretKey string
	TokenTTL  time.Duration
}

// MFA contains data for multi factor authentication.
type MFA struct {
	TemporaryTokenExpiration time.Duration
	AccessTokenExpiration    time.Duration
	RefreshTokenExpiration   time.Duration
}

// Email is configuration for email service.
type Email struct {
	SenderEmail              string
	APIKeyPublic             string
	APIKeyPrivate            string
	ForgotPasswordTemplateID int
}
