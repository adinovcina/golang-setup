package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/adinovcina/golang-setup/tools/logger"
)

type EnvironmentVariable string

func (e EnvironmentVariable) String() string {
	return string(e)
}

func Get(e EnvironmentVariable) string {
	return os.Getenv(e.String())
}

func MustGet(e EnvironmentVariable) string {
	v := os.Getenv(e.String())
	if v == "" {
		logger.Fatal().Msgf(" variable `%s` is not present in ENVIRONMENT", e.String())
	}

	return v
}

func GetOr(e EnvironmentVariable, fallback string) string {
	if value := Get(e); value != "" {
		return value
	}

	return fallback
}

func GetUint(e EnvironmentVariable) (uint, error) {
	intValue, err := GetInt(e)
	if err != nil {
		return 0, err
	}

	if intValue < 0 {
		return 0, errors.New("uint value must not be less than zero")
	}

	return uint(intValue), nil
}

func GetInt(e EnvironmentVariable) (int, error) {
	stringValue := Get(e)
	if stringValue == "" {
		// Not provided
		return 0, nil
	}

	result, err := strconv.Atoi(stringValue)
	if err != nil {
		return 0, fmt.Errorf("failed to parse value: %w", err)
	}

	return result, nil
}

func GetIntOr(e EnvironmentVariable, fallback int) int {
	stringValue := Get(e)
	if stringValue == "" {
		return fallback
	}

	result, err := strconv.Atoi(stringValue)
	if err != nil {
		logger.Fatal().Msgf("variable `%s` cannot be parsed to INTEGER", e.String())
	}

	return result
}

func GetBoolean(e EnvironmentVariable) (bool, error) {
	stringValue := Get(e)
	if stringValue == "" {
		// Not provided
		return false, nil
	}

	return strconv.ParseBool(stringValue)
}

func GetBooleanOr(e EnvironmentVariable, fallback bool) bool {
	val, err := strconv.ParseBool(Get(e))
	if err != nil {
		return fallback
	}

	return val
}

func GetDateTime(e EnvironmentVariable, fallback time.Duration) time.Duration {
	stringValue := Get(e)

	dateDuration, err := time.ParseDuration(stringValue)
	if err != nil {
		logger.Fatal().Msgf("variable `%s` cannot be parsed to DATE TIME", e.String())

		return fallback
	}

	return dateDuration
}
