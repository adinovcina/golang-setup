package api

import (
	"net/http"
	"strings"

	status "github.com/adinovcina/golang-setup/tools/network/statuscodes"
	"github.com/twinj/uuid"
)

// CreateAccountRequest used when admin sends request to the user to create new account for him.
type CreateAccountRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Validate request and decode into BaseResponse.
func (car *CreateAccountRequest) Validate(r *http.Request) (bool, *BaseResponse) {
	return ValidateRequestData(car, r, func() (bool, *BaseResponse) {
		response := new(BaseResponse)

		// Validate body params
		if strings.TrimSpace(car.Email) == "" {
			response.Error(status.ErrorMissingEmail)
		}

		if !validateEmail(car.Email) {
			response.Error(status.ErrorEmailNotInCorrectFormat)
		}

		return response.HasErrors(), response
	})
}

// ForgotPasswordRequest used when user wants to reset password.
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// Validate ForgotPasswordRequest.
func (fpr *ForgotPasswordRequest) Validate(r *http.Request) (bool, *BaseResponse) {
	return ValidateRequestData(fpr, r, func() (bool, *BaseResponse) {
		response := new(BaseResponse)

		// Validate body params
		if strings.TrimSpace(fpr.Email) == "" {
			response.Error(status.ErrorMissingEmail)
		}

		if !validateEmail(fpr.Email) {
			response.Error(status.ErrorEmailNotInCorrectFormat)
		}

		return response.HasErrors(), response
	})
}

// SetPasswordRequest used when user want to set a new password.
type SetPasswordRequest struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

// Validate SetPasswordRequest.
func (spr *SetPasswordRequest) Validate(r *http.Request) (bool, *BaseResponse) {
	return ValidateRequestData(spr, r, func() (bool, *BaseResponse) {
		response := new(BaseResponse)

		// Validate body params
		if strings.TrimSpace(spr.Password) == "" {
			response.Error(status.ErrorMissingPassword)
		}

		if strings.TrimSpace(spr.Token) == "" {
			response.Error(status.ErrorMissingToken)
		}

		return response.HasErrors(), response
	})
}

// UserActivateRequest contains id of user that needs to be deactivated / activated.
type UserActivateRequest struct {
	UserID uuid.UUID `json:"userID"`
}

// Validate UserActivateRequest.
func (uar *UserActivateRequest) Validate(r *http.Request) (bool, *BaseResponse) {
	return ValidateRequestData(uar, r, func() (bool, *BaseResponse) {
		response := new(BaseResponse)

		// Validate body params
		if _, err := uuid.Parse(uar.UserID.String()); err != nil {
			response.Error(status.ErrorMissingUserID)
		}

		return response.HasErrors(), response
	})
}

// ChangePasswordRequest contains new and old password for user to change.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// Validate ChangePasswordRequest.
func (cpr *ChangePasswordRequest) Validate(r *http.Request) (bool, *BaseResponse) {
	return ValidateRequestData(cpr, r, func() (bool, *BaseResponse) {
		response := new(BaseResponse)

		// Validate body params
		if strings.TrimSpace(cpr.NewPassword) == "" {
			response.Error(status.ErrorMissingPassword)
		}

		if strings.TrimSpace(cpr.CurrentPassword) == "" {
			response.Error(status.ErrorMissingPassword)
		}

		return response.HasErrors(), response
	})
}

// LogoutRequest used to logout user from our platform.
type LogoutRequest struct {
	Token string `json:"token"`
}

// Validate LogoutRequest.
func (lr *LogoutRequest) Validate(r *http.Request) (bool, *BaseResponse) {
	return ValidateRequestData(lr, r, func() (bool, *BaseResponse) {
		response := new(BaseResponse)

		// Validate body params
		if lr.Token == "" {
			response.Error(status.ErrorMissingToken)
		}

		return response.HasErrors(), response
	})
}
