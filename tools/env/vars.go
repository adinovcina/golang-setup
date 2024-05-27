package env

const (
	// SERVICE ENV VARIABLES.
	ServicePort        EnvironmentVariable = "SERVICE_PORT"
	ServiceEnvironment EnvironmentVariable = "SERVICE_ENVIRONMENT"
	LogLevel           EnvironmentVariable = "LOG_LEVEL"

	// TTL ENV VARIABLES.
	MFATemporaryTokenExpiration EnvironmentVariable = "MFA_TEMPORARY_TOKEN_EXPIRATION"
	MFARefreshTokenExpiration   EnvironmentVariable = "MFA_REFRESH_TOKEN_EXPIRATION"
	MFAAccessTokenExpiration    EnvironmentVariable = "MFA_ACCESS_TOKEN_EXPIRATION"

	// ACCOUNT ENV VARIABLES.
	MaxLoginFailures EnvironmentVariable = "MAX_LOGIN_FAILURES"
	BanDurationTime  EnvironmentVariable = "BAN_DURATION_TIME"

	// DATABASE ENV VARIABLES.
	DatabaseUsername         EnvironmentVariable = "DATABASE_USERNAME"
	DatabasePassword         EnvironmentVariable = "DATABASE_PASSWORD"
	DatabaseName             EnvironmentVariable = "DATABASE_NAME"
	DatabaseAddress          EnvironmentVariable = "DATABASE_ADDRESS"
	DatabasePort             EnvironmentVariable = "DATABASE_PORT"
	DatabaseMigrationFolder  EnvironmentVariable = "DATABASE_MIGRATION_FOLDER"
	DatabaseMigrationEnabled EnvironmentVariable = "DATABASE_MIGRATION_ENABLED"

	// REDIS ENV VARIABLES.
	RedisAddress   EnvironmentVariable = "REDIS_ADDRESS"
	RedisDatabase  EnvironmentVariable = "REDIS_DATABASE"
	RedisPassword  EnvironmentVariable = "REDIS_PASSWORD"
	RedisSecretKey EnvironmentVariable = "REDIS_SECRET_KEY"
	RedisTokenTTL  EnvironmentVariable = "REDIS_TOKEN_TTL"

	// Email ENV VARIABLES.
	APIKeyPublic             EnvironmentVariable = "API_KEY_PUBLIC"
	APIKeyPrivate            EnvironmentVariable = "API_KEY_PRIVATE"
	ForgotPasswordTemplateID EnvironmentVariable = "FORGOT_PASSWORD_TEMPLATE_ID"
	SenderEmail              EnvironmentVariable = "SENDER_EMAIL"

	// ENCRYPTION ENV VARIABLES.
	EncryptionProviderHashKey EnvironmentVariable = "ENCRYPTION_PROVIDER_HASH_KEY"
)
