package api

import (
	"net/http"
	"strings"

	status "github.com/adinovcina/golang-setup/tools/network/statuscodes"
	"github.com/twinj/uuid"
)

// AuthorizeRequest used when user send request to authorize to the app.
type AuthorizeRequest struct {
	Token string `json:"token"`
}

// Validate AuthorizeRequest.
func (ar *AuthorizeRequest) Validate(r *http.Request) (bool, *BaseResponse) {
	return ValidateRequestData(ar, r, func() (bool, *BaseResponse) {
		response := new(BaseResponse)

		// Validate body params
		if strings.TrimSpace(ar.Token) == "" {
			response.Error(status.ErrorMissingToken)
		}

		return response.HasErrors(), response
	})
}

// RefreshTokenRequest used when user wants to refresh his token.
type RefreshTokenRequest struct {
	Token string `json:"token"`
}

// Validate RefreshTokenRequest.
func (rtr *RefreshTokenRequest) Validate(r *http.Request) (bool, *BaseResponse) {
	return ValidateRequestData(rtr, r, func() (bool, *BaseResponse) {
		response := new(BaseResponse)

		// Validate body params
		if strings.TrimSpace(rtr.Token) == "" {
			response.Error(status.ErrorMissingToken)
		}

		return response.HasErrors(), response
	})
}

// LoginDataResponse contains response data after login is called.
type LoginDataResponse struct {
	Token    Token     `json:"token"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	Language string    `json:"language"`
	UserID   uuid.UUID `json:"userID"`
}
