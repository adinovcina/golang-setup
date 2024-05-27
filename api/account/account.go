package account

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/adinovcina/golang-setup/api"
	m "github.com/adinovcina/golang-setup/api/middleware"
	"github.com/adinovcina/golang-setup/config"
	"github.com/adinovcina/golang-setup/store"
	"github.com/adinovcina/golang-setup/tools/encryption"
	"github.com/adinovcina/golang-setup/tools/logger"
	mailjet "github.com/adinovcina/golang-setup/tools/mailjet"
	status "github.com/adinovcina/golang-setup/tools/network/statuscodes"
	"github.com/adinovcina/golang-setup/tools/utils"

	"github.com/go-chi/chi/v5"
)

func AttachAccountRoutes(r chi.Router,
	conf *config.Config,
	repo store.Repository,
	inMemRepo store.InMemRepository,
	mailjetClient *mailjet.Client,
) {
	svc := newService(conf, repo, inMemRepo, mailjetClient)

	// Unprotected REST routes for "account" resource
	r.Route("/account", func(r chi.Router) {
		// Authenticate user using email and password
		r.Post("/authenticate", svc.handleAuthenticateUser)
		// Used by users to authorize their selected user role
		r.Post("/authorize", svc.handleAuthorizeUser)
		// Used by user to extend his session once it's expired
		r.Post("/refresh-token", svc.handleRefreshToken)
		// Used by user to reset their forgotten password
		r.Post("/forgot-password", svc.handleForgotPassword)
		// Used by user to set his new password once he receive reset link on email
		r.Post("/set-password", svc.handleSetPassword)

		r.Group(func(r chi.Router) {
			// Private API group
			r.Use(m.AuthorizeRequest(&conf.Redis, inMemRepo))

			// Used by logged in user to fetch his user roles
			r.Get("/roles", svc.handleGetRoles)
			// Used by user to change their password
			r.Post("/change-password", svc.handleChangePassword)
			// Used by user to update their profile
			r.Patch("/users/profile", svc.handleUpdateUserProfile)
			// Used to logout user from platform
			r.Post("/logout", svc.handleLogout)
			// Used to fetch user profile
			r.Get("/me", svc.handleGetProfile)

			r.Group(func(r chi.Router) {
				// Restrict only to admin role
				r.Use(m.CheckAllowedRoles(store.GetRoles().Admin))
				// Used by admin user to activate or deactivate another user
				r.Post("/activate", svc.handleActivateUser)
				// Used by admin to retrieve list of all users in the system
				r.With(m.PaginationCursor(repo)).Get("/users", svc.handleGetUsers)
			})
		})
	})
}

// authenticateUser is used to authenticate user, and type of MFA user is using.
func (s *service) handleAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	requestData := api.RequestData(r)
	request := new(api.AuthenticateUserRequest)

	valid, response := request.Validate(r)
	response.RequestID = requestData.RequestID

	if !valid {
		logger.Error().Msgf("invalid request received: %v", response)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, nil)

		return
	}

	// Retrieve user from the database by email and check if exists
	user, err := s.repo.GetUserByEmail(request.Email)

	// If user is not found then tell user that email or password is incorrect
	if err != nil && err.Error() == store.UserNotFound {
		response.Error(status.ErrorIncorrectEmailOrPassword)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	} else if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	// Handle case when user is not active
	if !user.Active {
		response.Error(status.ErrorUserNotActive)
		api.ErrorResponse(response, http.StatusUnauthorized, w, r, err)

		return
	}

	if user.VerifyIfUserIsSuspended(s.conf.Account.MaxLoginFailures) {
		response.Errors = append(response.Errors,
			api.Error{Code: status.ErrorUserSuspended, Message: fmt.Sprintf("%s %s", status.ErrorStatusText(status.ErrorUserSuspended), user.LoginBlockedUntil)})
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	}

	// All good so far, now check if password match. Password is hashed in database
	// NOTE: We will not show user that he missed his password since that would be easy for
	// hackers to guess that email is correct.
	err = encryption.IsValid(user.Password, request.Password)
	if err != nil {
		response.Error(status.ErrorIncorrectEmailOrPassword)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		// We want to increment the failed login counter every time a user misses his password.
		_, err = s.repo.UpdateLoginAttempt(user.ID, s.conf.Account.BanDurationTime.Minutes(), s.conf.Account.MaxLoginFailures)
		if err != nil {
			api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)
		}

		return
	}

	// Password is valid, next step is to create a token
	temporaryToken, err := s.createTemporaryToken(user)
	if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	response.Data = api.AuthenticateUserDataResponse{
		Token: temporaryToken,
	}

	api.SuccessResponse(response, http.StatusOK, w)
}

// handleAuthorizeUser is used to authorize user, create user's session and return token and refresh token.
func (s *service) handleAuthorizeUser(w http.ResponseWriter, r *http.Request) {
	requestData := api.RequestData(r)
	request := new(api.AuthorizeRequest)

	// Validate request sent from the frontend
	valid, response := request.Validate(r)
	response.RequestID = requestData.RequestID

	if !valid {
		logger.Error().Msgf("invalid request received: %v", request)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, nil)

		return
	}

	// Retrieve user from the database by token and check if exists
	user, err := s.repo.GetUserByToken(request.Token, store.GetTokenTypes().MFA)

	// If user is not found then tell user that email or password is incorrect
	if err != nil && err.Error() == store.UserNotFound {
		response.Error(status.ErrorIncorrectEmailOrPassword)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	} else if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	// If temporary token has expired return unauthorized
	if user.Expired {
		response.Error(status.ErrorTokenExpiredOrNotValid)
		api.ErrorResponse(response, http.StatusUnauthorized, w, r, err)

		return
	}

	// Handle case when user is not active
	if !user.Active {
		response.Error(status.ErrorUserNotActive)
		api.ErrorResponse(response, http.StatusUnauthorized, w, r, err)

		return
	}

	// Reset the login counter to zero when the user has successfully logged in
	if user.FailedLoginCount > 0 {
		err = s.repo.ResetFailedLoginCounter(user.ID)
		if err != nil {
			api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

			return
		}
	}

	token, err := s.createToken(r.Context(), user)
	if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	refreshToken, err := s.createRefreshToken(user.ID)
	if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	response.Data = api.LoginDataResponse{
		Email: user.Email,
		Name:  user.Name,
		Token: api.Token{
			Token:        token,
			RefreshToken: refreshToken,
		},
		Role:     user.Role,
		UserID:   user.ID,
		Language: user.Language,
	}

	api.SuccessResponse(response, http.StatusOK, w)
}

// handleGetRoles is used retrieve user roles associated to him.
func (s *service) handleGetRoles(w http.ResponseWriter, r *http.Request) {
	requestData := api.RequestData(r)
	response := new(api.BaseResponse)
	response.RequestID = requestData.RequestID

	roles, err := s.repo.GetUserRoles(requestData.UserID)
	if err != nil {
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)
		return
	}

	response.Data = roles

	api.SuccessResponse(response, http.StatusOK, w)
}

// handleForgotPassword used by user to send an email with password reset link.
func (s *service) handleForgotPassword(w http.ResponseWriter, r *http.Request) {
	requestData := api.RequestData(r)
	request := new(api.ForgotPasswordRequest)

	valid, response := request.Validate(r)
	response.RequestID = requestData.RequestID

	if !valid {
		logger.Error().Msgf("invalid request received: %v", request)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, nil)

		return
	}

	user, err := s.repo.GetUserByEmail(request.Email)
	if err != nil {
		logger.Error().Err(err).Msgf("ForgotPassword unable to find an account: %v.", request.Email)

		response.Error(status.ErrorEmailDoesNotExists)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	}

	if !user.Active {
		logger.Error().Err(err).Msgf("ForgotPassword user is not active email: %v and user id: %v.",
			request.Email, user.ID)

		response.Error(status.ErrorUserNotActive)
		api.ErrorResponse(response, http.StatusUnauthorized, w, r, err)

		return
	}

	// Generate password token. For create an account this link sent to email should be active for 30 days
	token := strings.ReplaceAll(utils.GenerateUniqueID()+utils.GenerateUniqueID()+utils.GenerateUniqueID(), "-", "")
	tokenExpiresAt := time.Now().Add(s.conf.MFA.AccessTokenExpiration).Unix()

	passwordToken, err := s.repo.AddPasswordResetToken(user.ID, token, tokenExpiresAt)
	if err != nil {
		logger.Error().Err(err).Msgf(`ForgotPassword unable to create password token code for email: %v and user id: %v.`,
			request.Email, user.ID)
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	// Send email in a new thread
	go s.mailjetClient.SendEmailResetPassword(s.conf.Email.ForgotPasswordTemplateID, user.Email,
		s.conf.Email.SenderEmail, user.Email, passwordToken.Token)

	api.SuccessResponse(response, http.StatusNoContent, w)
}

// handleSetPassword will set password if token is valid.
func (s *service) handleSetPassword(w http.ResponseWriter, r *http.Request) {
	requestData := api.RequestData(r)
	request := new(api.SetPasswordRequest)

	valid, response := request.Validate(r)
	response.RequestID = requestData.RequestID

	if !valid {
		logger.Error().Msgf("invalid request received: %v", request)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, nil)

		return
	}

	// Retrieve user id by token from ResetToken
	token, err := s.repo.GetPasswordTokenByToken(request.Token)
	if err != nil || (token != nil && token.ExpiresAt < time.Now().Unix()) {
		response.Error(status.ErrorTokenExpiredOrNotValid)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	}

	password, err := encryption.Encrypt(request.Password)
	if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	user, err := s.repo.SetPassword(token.UserID, password, request.Token)
	if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	sessionToken, err := s.createToken(r.Context(), user)
	if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)
		return
	}

	refreshToken, err := s.createRefreshToken(user.ID)
	if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)
		return
	}

	response.Data = api.LoginDataResponse{
		UserID:   user.ID,
		Language: user.Language,
		Email:    user.Email,
		Name:     user.Name,
		Token: api.Token{
			Token:        sessionToken,
			RefreshToken: refreshToken,
		},
		Role: user.Role,
	}

	api.SuccessResponse(response, http.StatusOK, w)
}

// handleGetUsers will retrieve list of the users.
func (s *service) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	requestData := api.RequestData(r)
	filter := &store.UserFilter{}

	response := &api.BaseResponse{}
	response.RequestID = requestData.RequestID

	activeStr := r.URL.Query().Get("active")
	if activeStr != "" {
		userActive, err := strconv.ParseBool(activeStr)
		if err != nil {
			response.Error(status.ErrorInvalidQueryURLParameters)
			api.ErrorResponse(response, http.StatusNotFound, w, r, err)

			return
		}

		filter.Active = &userActive
	}

	search := r.URL.Query().Get("search")
	if search != "" {
		filter.Search = &search
	}

	users, err := s.repo.GetUsers(filter)
	if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)
		return
	}

	response.Data = api.PaginatedCursorResponse{
		Results:    users,
		Pagination: s.repo.PaginatorCursor(),
	}

	api.SuccessResponse(response, http.StatusOK, w)
}

func (s *service) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	requestData := api.RequestData(r)
	request := new(api.RefreshTokenRequest)

	valid, response := request.Validate(r)
	response.RequestID = requestData.RequestID

	if !valid {
		logger.Error().Msgf("invalid request received: %v", request)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, nil)

		return
	}

	// Get refresh loginToken from DB
	loginToken, err := s.repo.GetTokenByTokenAndType(request.Token, store.GetTokenTypes().RefreshToken)
	if err != nil {
		response.Error(status.ErrorMissingToken)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	}

	// Fetch user by ID
	user, err := s.repo.GetUserByID(loginToken.UserID)
	if err != nil {
		response.Error(status.ErrorGetUser)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	}

	// Validate token expiration
	if loginToken.Expired {
		response.Error(status.ErrorTokenExpiredOrNotValid)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	}

	// Delete old refresh token
	if err = s.repo.DeleteTokenByID(loginToken.ID); err != nil {
		response.Error(status.ErrorDeleteToken)
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	// Generate access token
	token, err := s.createToken(r.Context(), user)
	if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)
		return
	}

	// Generate refresh token
	refreshToken, err := s.createRefreshToken(user.ID)
	if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)
		return
	}

	// Return final response data
	response.Data = api.LoginDataResponse{
		Email: user.Email,
		Name:  user.Name,
		Token: api.Token{
			Token:        token,
			RefreshToken: refreshToken,
		},
		Role:     user.Role,
		UserID:   user.ID,
		Language: user.Language,
	}

	api.SuccessResponse(response, http.StatusOK, w)
}

// Activate will activate or deactivate user.
func (s *service) handleActivateUser(w http.ResponseWriter, r *http.Request) {
	requestData := api.RequestData(r)
	request := new(api.UserActivateRequest)

	valid, response := request.Validate(r)
	response.RequestID = requestData.RequestID

	if !valid {
		logger.Error().Msgf("invalid request received: %v", request)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, nil)

		return
	}

	if request.UserID == requestData.UserID {
		response.Error(status.ErrorUnableToDeactivateAdmin)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, nil)

		return
	}

	err := s.repo.ActivateUser(request.UserID)
	if err != nil {
		response.Error(status.ErrorActivateUser)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	}

	api.SuccessResponse(response, http.StatusOK, w)
}

// handleChangePassword changes password for currently logged user.
func (s *service) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	requestData := api.RequestData(r)
	request := new(api.ChangePasswordRequest)

	valid, response := request.Validate(r)
	response.RequestID = requestData.RequestID

	if !valid {
		logger.Error().Msgf("invalid request received: %v", request)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, nil)

		return
	}

	// Retrieve logged user
	user, err := s.repo.GetUserByID(requestData.UserID)
	if err != nil {
		response.Error(status.ErrorGetUser)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	}

	// Check if current supplied password match password in DB
	err = encryption.IsValid(user.Password, request.CurrentPassword)
	if err != nil {
		response.Error(status.ErrorCurrentPasswordMismatch)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	}

	// All good save new password into database
	newPassword, err := encryption.Encrypt(request.NewPassword)
	if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	err = s.repo.SetNewPassword(requestData.UserID, newPassword)
	if err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	api.SuccessResponse(response, http.StatusOK, w)
}

// handleUpdateUserProfile updates user profile.
func (s *service) handleUpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	requestData := api.RequestData(r)
	request := new(api.UpdateUserProfileRequest)

	valid, response := request.Validate(r)
	response.RequestID = requestData.RequestID

	if !valid {
		logger.Error().Msgf("invalid request received: %v", request)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, nil)

		return
	}

	if requestData.UserID != request.UserID && requestData.Role != store.GetRoles().Admin.Name {
		api.ErrorResponse(response, http.StatusForbidden, w, r, nil)

		return
	}

	user, err := s.repo.GetUserByID(request.UserID)
	if err != nil {
		response.Error(status.ErrorGetUser)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	}

	if request.Name != nil && *request.Name != user.Name {
		user.Name = *request.Name
	}

	if request.Phone != nil && *request.Phone != user.Phone {
		user.Phone = *request.Phone
	}

	user, err = s.repo.UpdateUser(user)
	if err != nil {
		logger.Error().Err(err).Msgf("unable to update user with id %v", user.ID)
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	response.Data = api.UserProfileDataResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Language: user.Language,
		Role:     user.Role,
	}

	api.SuccessResponse(response, http.StatusOK, w)
}

// handleLogout logs the user out of the application and deletes all their sessions.
func (s *service) handleLogout(w http.ResponseWriter, r *http.Request) {
	request := new(api.LogoutRequest)

	valid, response := request.Validate(r)
	if !valid {
		logger.Error().Msgf("invalid request received: %v", request)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, nil)

		return
	}

	ctx := r.Context()

	middlewareData := api.MiddlewareDataFromContext(ctx)
	if middlewareData != nil && middlewareData.SessionKey != "" {
		if err := s.inMemRepo.DelSessionWithKey(ctx, middlewareData.SessionKey); err != nil {
			logger.Warn().Err(err).Msgf("failed to delete session with key %v", middlewareData.SessionKey)
		}
	}

	// Retrieve token from the database by token and type
	token, err := s.repo.GetTokenByTokenAndType(request.Token, store.GetTokenTypes().RefreshToken)
	if err != nil {
		response.Error(status.ErrorMissingToken)
		api.ErrorResponse(response, http.StatusNotFound, w, r, err)

		return
	}

	// Delete refresh token
	if err := s.repo.DeleteTokenByID(token.ID); err != nil {
		api.ErrorResponse(response, http.StatusInternalServerError, w, r, err)

		return
	}

	api.SuccessResponse(response, http.StatusOK, w)
}

// handleGetProfile gets user profile.
func (s *service) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	requestData := api.RequestData(r)
	response := &api.BaseResponse{}
	response.RequestID = requestData.RequestID

	user, err := s.repo.GetUserByID(requestData.UserID)
	if err != nil {
		response.Error(status.ErrorGetUser)
		api.ErrorResponse(response, http.StatusBadRequest, w, r, err)

		return
	}

	response.Data = user

	api.SuccessResponse(response, http.StatusOK, w)
}
