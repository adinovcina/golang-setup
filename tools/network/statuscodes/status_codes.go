package status

// Code - used for our custom error codes to specify for the forntend what exactly happened.
type Code int

const (
	// InternalServerError - error when something happens internally.
	InternalServerError = 1000
	// IncorrectBodyFormat - error when decoded body is not in correct format.
	IncorrectBodyFormat = 1001
	// InvalidResponseBody - error when response body is not in a valid format.
	InvalidResponseBody = 1002
	// EmptyBody - error when body is not supplied.
	EmptyBody = 1003
	// ErrorMissingToken - error when token is not sent or when it all empty spaces.
	ErrorMissingToken = 1004
	// ErrorMissingEmail - error when email is not sent or when it all empty spaces.
	ErrorMissingEmail = 1005
	// ErrorMissingPassword - error when password is not sent.
	ErrorMissingPassword = 1006
	// ErrorEmailOrPasswordNotMatch - error when password or email do not match.
	ErrorEmailOrPasswordNotMatch = 1007
	// ErrorEmailNotInCorrectFormat - error when email is not in correct format.
	ErrorEmailNotInCorrectFormat = 1008
	// ErrorIncorrectEmailOrPassword used by service to indicate that user missed login credentials.
	ErrorIncorrectEmailOrPassword = 1009
	// ErrorUserNotActive used when user is not active due to multiple failed login attempts or stil not activated via email.
	ErrorUserNotActive = 1010
	// ErrorTokenExpiredOrNotValid used when token expired or not valid.
	ErrorTokenExpiredOrNotValid = 1011
	// ErrorMissingPhone - error when code is not sent or when it all empty spaces.
	ErrorMissingPhone = 1012
	// ErrorInvalidQueryURLParameters - error when decoded url query params is not in correct format.
	ErrorInvalidQueryURLParameters = 1013
	// ErrorActivateUser used when unable to activate user.
	ErrorActivateUser = 1014
	// ErrorMissingUserID used when user id is not supplied.
	ErrorMissingUserID = 1015
	// ErrorUnableToDeactivateAdmin used when admin wants to deactivate / activate himself.
	ErrorUnableToDeactivateAdmin = 1016
	// ErrorCurrentPasswordMismatch used when password and stored password mismatch.
	ErrorCurrentPasswordMismatch = 1017
	// ErrorDeleteToken is used when error happens during the delete token operation.
	ErrorDeleteToken = 1018
	// ErrorGetUser is used when unable to get user.
	ErrorGetUser = 1019
	// ErrorEmailDoesNotExists is used when email email does not exists.
	ErrorEmailDoesNotExists = 1020
	// ErrorUserSuspended used when user is suspended due to multiple failed login attempts.
	ErrorUserSuspended = 1021
)

// / ****************************************************
// / INFO
// / ****************************************************.
func errToStatusTextMap() map[int]string { //nolint:funlen // ignore
	statusText := map[int]string{
		InternalServerError:            "internal server error",
		EmptyBody:                      "empty request body",
		IncorrectBodyFormat:            "incorrect body format",
		InvalidResponseBody:            "invalid response body",
		ErrorMissingEmail:              "missing parameter email",
		ErrorMissingPassword:           "missing parameter password",
		ErrorEmailOrPasswordNotMatch:   "email or password do not match",
		ErrorEmailNotInCorrectFormat:   "email not in correct format",
		ErrorIncorrectEmailOrPassword:  "incorrect email or password",
		ErrorUserNotActive:             "user is not active",
		ErrorTokenExpiredOrNotValid:    "token expired or not valid",
		ErrorMissingPhone:              "missing phone paramater",
		ErrorInvalidQueryURLParameters: "incorrect URL query params format",
		ErrorMissingUserID:             "missing user id",
		ErrorUnableToDeactivateAdmin:   "unable to deactivate or activate admin account",
		ErrorActivateUser:              "unable to activate account",
		ErrorDeleteToken:               "unable to delete token",
		ErrorCurrentPasswordMismatch:   "current password mismatch",
		ErrorGetUser:                   "unable to fetch user",
		ErrorMissingToken:              "missing or invalid token",
		ErrorEmailDoesNotExists:        "email does not exists",
		ErrorUserSuspended:             "user suspended until",
	}

	return statusText
}

// ErrorStatusText returns the associated status text for error code.
func ErrorStatusText(code int) string {
	if v, ok := errToStatusTextMap()[code]; ok {
		return v
	}

	return "invalid request"
}
