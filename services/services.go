package services

import (
	"github.com/adinovcina/golang-setup/config"
	"github.com/adinovcina/golang-setup/tools/mailjet"
)

type AppServices struct {
	mailjetService *mailjet.Client
}

// Init will initialize services.
func Init(appConfig *config.Config) *AppServices {
	// Initialize Mailjet
	mailjetService := newMailjetService(appConfig.Email)

	return &AppServices{
		mailjetService: mailjetService,
	}
}

// GetMailjetClient returns the Mailjet client.
func (s *AppServices) GetMailjetClient() *mailjet.Client {
	return s.mailjetService
}
