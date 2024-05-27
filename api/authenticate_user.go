package api

import (
	"net/http"
	"regexp"
	"strings"

	status "github.com/adinovcina/golang-setup/tools/network/statuscodes"
)

// AuthenticateUserRequest used when user send request to login to the app.
type AuthenticateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate AuthenticateUserRequest.
func (aur *AuthenticateUserRequest) Validate(r *http.Request) (bool, *BaseResponse) {
	return ValidateRequestData(aur, r, func() (bool, *BaseResponse) {
		response := new(BaseResponse)

		if strings.TrimSpace(aur.Email) == "" {
			response.Error(status.ErrorMissingEmail)
		}

		if strings.TrimSpace(aur.Password) == "" {
			response.Error(status.ErrorMissingPassword)
		}

		if !validateEmail(aur.Email) {
			response.Error(status.ErrorEmailNotInCorrectFormat)
		}

		return response.HasErrors(), response
	})
}

// AuthenticateUserDataResponse contains response data after login is called.
type AuthenticateUserDataResponse struct {
	Token string `json:"token"`
}

// ValidateEmail will check if email match to regex.
func validateEmail(email string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%\-+]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`).MatchString(email)
}
