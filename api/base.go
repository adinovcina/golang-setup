package api

import (
	"encoding/json"
	"net/http"

	"github.com/adinovcina/golang-setup/tools/logger"
	status "github.com/adinovcina/golang-setup/tools/network/statuscodes"
	"github.com/adinovcina/golang-setup/tools/paging"
)

// BaseResponse structure that will have on all responses from all APIs.
type BaseResponse struct {
	Data      interface{} `json:"data"`
	RequestID string      `json:"requestID"`
	Errors    []Error     `json:"errors"`
}

// PaginatedResponse contains response when we use basic pagination.
type PaginatedResponse struct {
	Results    any              `json:"results"`
	Pagination paging.Paginator `json:"pagination"`
}

// PaginatedCursorResponse contains response when we use cursos pagination.
type PaginatedCursorResponse struct {
	Results    any                    `json:"results"`
	Pagination paging.PaginatorCursor `json:"pagination"`
}

// Error object that will contain details of the error.
type Error struct {
	Message string `json:"errorMessage"`
	Code    int    `json:"errorCode"`
}

// Token contains all important data for the token.
type Token struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func (l *BaseResponse) Error(customStatus int) {
	l.Errors = append(l.Errors, Error{Code: customStatus, Message: status.ErrorStatusText(customStatus)})
}

func (r *BaseResponse) HasErrors() bool {
	return len(r.Errors) <= 0
}

func ErrorResponse(response *BaseResponse, statusCode int, w http.ResponseWriter, r *http.Request, err error) {
	// If error is not nil, log the error
	if err != nil {
		requestData := RequestData(r)
		logger.Error().Err(err).Msgf("request_id: %s", requestData.RequestID)
	}

	// If there is no explicit error defined, return internal error message
	if len(response.Errors) == 0 {
		response.Error(status.InternalServerError)
	}

	jsonResponse(response, statusCode, w)
}

// RequestData to retrieve the Data object from the context.
// The *Data* object will be initialized if not present already.
func RequestData(r *http.Request) *Data {
	return MiddlewareDataFromContext(r.Context())
}

func SuccessResponse(response *BaseResponse, statusCode int, w http.ResponseWriter) {
	if response == nil {
		response = new(BaseResponse)
	}

	jsonResponse(response, statusCode, w)
}

func jsonResponse(response *BaseResponse, statusCode int, w http.ResponseWriter) {
	// Set content type only if there is response data
	if statusCode != http.StatusNoContent {
		w.Header().Set("Content-Type", "application/json")
	}

	marshaledData, err := json.Marshal(response)
	if err != nil {
		logger.Error().Err(err).Msg("marshaling response data failed")

		response = &BaseResponse{
			Errors: []Error{{Code: status.InvalidResponseBody, Message: status.ErrorStatusText(status.InvalidResponseBody)}},
		}

		marshaledData, _ = json.Marshal(response)
	}

	w.WriteHeader(statusCode)
	writeResponseData(statusCode, marshaledData, w)
}

func writeResponseData(statusCode int, data []byte, w http.ResponseWriter) {
	if statusCode != http.StatusNoContent {
		if _, err := w.Write(data); err != nil {
			logger.Error().Err(err).Msg("write failed")
		}
	}
}

func ValidateRequestData(requestData interface{},
	r *http.Request,
	validationFunc func() (bool, *BaseResponse),
) (bool, *BaseResponse) {
	response := new(BaseResponse)

	// Validate if body is empty since we require some input
	if r.Body == nil {
		response.Error(status.EmptyBody)

		return false, response
	}

	// Decode body
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		response.Error(status.IncorrectBodyFormat)

		return false, response
	}

	defer r.Body.Close()

	return validationFunc()
}
