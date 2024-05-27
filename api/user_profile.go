package api

import (
	"net/http"
	"strings"

	status "github.com/adinovcina/golang-setup/tools/network/statuscodes"
	"github.com/twinj/uuid"
)

// UpdateUserProfileRequest contains user profile info.
type UpdateUserProfileRequest struct {
	Name   *string   `json:"name"`
	Phone  *string   `json:"phone"`
	UserID uuid.UUID `json:"userID"`
}

// Validate UpdateUserProfileRequest.
func (uupr *UpdateUserProfileRequest) Validate(r *http.Request) (bool, *BaseResponse) {
	return ValidateRequestData(uupr, r, func() (bool, *BaseResponse) {
		response := new(BaseResponse)

		// Validate body params
		if uupr.Phone != nil && strings.TrimSpace(*uupr.Phone) == "" {
			response.Error(status.ErrorMissingPhone)
		}

		if _, err := uuid.Parse(uupr.UserID.String()); err != nil {
			response.Error(status.ErrorMissingUserID)
		}

		return response.HasErrors(), response
	})
}

type UserProfileDataResponse struct {
	Name     string    `json:"name,omitempty"`
	Email    string    `json:"email,omitempty"`
	Phone    string    `json:"phone,omitempty"`
	Language string    `json:"language,omitempty"`
	Role     string    `json:"role,omitempty"`
	ID       uuid.UUID `json:"id,omitempty"`
}
