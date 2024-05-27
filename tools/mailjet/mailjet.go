package mailjet

import (
	"github.com/adinovcina/golang-setup/tools/logger"
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type Client struct {
	client *mailjet.Client
	apiKey string
	domain string
}

// NewClient email service.
func NewClient(apiKeyPublic, apiKeyPrivate string) *Client {
	client := mailjet.NewMailjetClient(apiKeyPublic, apiKeyPrivate)

	return &Client{
		apiKey: apiKeyPublic,
		domain: apiKeyPrivate,
		client: client,
	}
}

// SendEmailResetPassword will send email to user to reset his password.
func (c *Client) SendEmailResetPassword(templateID int, name, fromEmail, toEmail, token string) {
	// Define the variables for the template
	vars := map[string]interface{}{
		"mj_reset_password_link": "https://example.com/reset-password/" + token,
		"mj_user_name":           name,
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: fromEmail,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: toEmail,
				},
			},
			TemplateID:       templateID,
			TemplateLanguage: true,
			Variables:        vars,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}

	// Send the email
	_, err := c.client.SendMailV31(&messages)
	if err != nil {
		logger.Error().Msgf("Mailjet responded with error: %v", err)
	}
}
