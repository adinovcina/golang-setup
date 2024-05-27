package services

import (
	"github.com/adinovcina/golang-setup/config"
	"github.com/adinovcina/golang-setup/tools/mailjet"
)

// newMailjetService Initialize.
func newMailjetService(appConfig config.Email) *mailjet.Client {
	// INITIALIZE EMAIL SERVICE
	return mailjet.NewClient(
		appConfig.APIKeyPublic,
		appConfig.APIKeyPrivate,
	)
}
