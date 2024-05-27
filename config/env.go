package config

import (
	"github.com/adinovcina/golang-setup/tools/env"
)

func loadFromEnv() (*Config, error) {
	config := &Config{
		Service: Service{
			Port:        env.GetOr(env.ServicePort, apiPortDefault),
			Environment: env.GetOr(env.ServiceEnvironment, stageDevelopment),
			LogLevel:    env.GetOr(env.LogLevel, logLevelInfo),
		},
		Account: Account{
			MaxLoginFailures: env.GetIntOr(env.MaxLoginFailures, maxLoginFailures),
			BanDurationTime:  env.GetDateTime(env.BanDurationTime, banDurationDefaultTime),
		},
		Database: Database{
			Username:         env.MustGet(env.DatabaseUsername),
			Password:         env.MustGet(env.DatabasePassword),
			Name:             env.MustGet(env.DatabaseName),
			Address:          env.MustGet(env.DatabaseAddress),
			Port:             env.MustGet(env.DatabasePort),
			MigrationFolder:  env.MustGet(env.DatabaseMigrationFolder),
			MigrationEnabled: env.GetBooleanOr(env.DatabaseMigrationEnabled, migrationEnabledDefault),
		},
		Timeouts: Timeouts{
			ReadDuration:  timeoutDuration,
			WriteDuration: timeoutDuration,
		},
		MFA: MFA{
			TemporaryTokenExpiration: env.GetDateTime(env.MFATemporaryTokenExpiration, mfaTemporaryTokenExpirationDefault),
			AccessTokenExpiration:    env.GetDateTime(env.MFAAccessTokenExpiration, mfaAccessTokenExpirationDefault),
			RefreshTokenExpiration:   env.GetDateTime(env.MFARefreshTokenExpiration, mfaRefreshTokenExpirationDefault),
		},
		Redis: Redis{
			Address:   env.MustGet(env.RedisAddress),
			Database:  env.MustGet(env.RedisDatabase),
			Password:  env.Get(env.RedisPassword),
			SecretKey: env.MustGet(env.RedisSecretKey),
			TokenTTL:  env.GetDateTime(env.RedisTokenTTL, redisTokenExpirationDefault),
		},
		Email: Email{
			APIKeyPublic:             env.MustGet(env.APIKeyPublic),
			APIKeyPrivate:            env.MustGet(env.APIKeyPrivate),
			SenderEmail:              env.MustGet(env.SenderEmail),
			ForgotPasswordTemplateID: env.GetIntOr(env.ForgotPasswordTemplateID, 0),
		},
	}

	return config, nil
}
