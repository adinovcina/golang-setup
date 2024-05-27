package mysqlstore

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/adinovcina/golang-setup/store"
	"github.com/adinovcina/golang-setup/tools/logger"
	"github.com/twinj/uuid"
)

// Reset failed login counter to 0 after successful login.
func (r *Repository) ResetFailedLoginCounter(userID uuid.UUID) error {
	query, err := r.db.Prepare("CALL ResetFailedLoginCount(?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to PREPARE statement for: CALL ResetFailedLoginCount(%v). ", userID)
		return err
	}

	defer query.Close()

	_, err = query.Exec(userID)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to EXECUTE statement for:  CALL ResetFailedLoginCount(%v). ", userID)
		return err
	}

	return nil
}

// UpdateLoginAttempt used for saving failed login count when code is incorrect.
func (r *Repository) UpdateLoginAttempt(loggedUserID uuid.UUID, minutes float64, maxLoginFailures int) (int64, error) {
	var failedLoginCount int64

	query, err := r.db.Prepare("CALL UpdateLoginAttempt(?, ?, ?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL UpdateLoginAttempt(%v, %v, %v).",
			loggedUserID, minutes, maxLoginFailures)
		return failedLoginCount, err
	}

	defer query.Close()

	err = query.QueryRow(loggedUserID, minutes, maxLoginFailures).Scan(&failedLoginCount)
	if errors.Is(err, sql.ErrNoRows) {
		return failedLoginCount, nil
	}

	if err != nil {
		logger.Error().Err(err).Msgf("failed to execute statement: CALL UpdateLoginAttempt(%v, %v, %v).",
			loggedUserID, minutes, maxLoginFailures)
		return failedLoginCount, err
	}

	return failedLoginCount, nil
}

// AddLoginToken will add login token to DB and return roles associated with user
func (r *Repository) AddLoginToken(userID uuid.UUID, expirationTime int64, token, tokenType string) error {
	query, err := r.db.Prepare("CALL AddLoginToken(?, ?, ?, ?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL AddLoginToken(%v, %v, %v, %v).",
			userID, token, tokenType, expirationTime)
		return err
	}

	defer query.Close()

	rows, err := query.Query(userID, token, tokenType, expirationTime)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if err != nil {
		logger.Error().Err(err).Msgf("There was an error executing query: CALL AddLoginToken(%v, %v, %v, %v).",
			userID, token, tokenType, expirationTime)
		return err
	}

	defer rows.Close()

	return nil
}

// GetUserRoles will fetch all roles associated with user
func (r *Repository) GetUserRoles(userID uuid.UUID) ([]*store.Role, error) {
	query, err := r.db.Prepare("CALL GetRoles(?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL GetRoles(%v).", userID)
		return nil, err
	}

	defer query.Close()

	userRoles := make([]*store.Role, 0)

	rows, err := query.Query(userID)
	if errors.Is(err, sql.ErrNoRows) {
		return userRoles, nil
	}

	if err != nil {
		logger.Error().Err(err).Msgf("There was an error executing query: CALL GetRoles(%v).", userID)
		return nil, err
	}

	defer rows.Close()

	userRoles, err = fetchUserRoles(rows)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch login roles")
		return nil, err
	}

	return userRoles, nil
}

// SetPassword used by user to update his password field in database.
func (r *Repository) SetPassword(userID uuid.UUID, password, token string) (*store.User, error) {
	query, err := r.db.Prepare("CALL SetPassword(?, ?, ?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL SetPassword(%v, %v, %v).",
			userID, password, token)
		return nil, err
	}

	defer query.Close()

	user := new(store.User)

	err = query.QueryRow(userID, password, token).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Language, &user.Active,
			&user.Phone, &user.Password, &user.Role, &user.RoleID, &user.CreatedAt)
	if err != nil {
		logger.Error().Err(err).Msgf("There was an error executing query: CALL SetPassword(%v, %v, %v).",
			userID, password, token)
		return nil, err
	}

	return user, nil
}

// Activate will update user's active status.
func (r *Repository) ActivateUser(userID uuid.UUID) error {
	query, err := r.db.Prepare("CALL ActivateUser(?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL ActivateUser(%v)", userID)

		return err
	}

	defer query.Close()

	res, err := query.Exec(userID)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to execute statement CALL ActivateUser(%v)", userID)
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		logger.Error().Err(err).Msgf(
			"failed to execute statement CALL ActivateUser(%v)", userID)
		return err
	}

	if ra == 0 {
		err = fmt.Errorf("unable to activate user %v", userID)
		logger.Error().Err(err).Msgf("Rows Affected = 0 for statement: CALL ActivateUser(%v)", userID)

		return err
	}

	return nil
}

// SetNewPassword sets a new password for logged user.
func (r *Repository) SetNewPassword(userID uuid.UUID, password string) error {
	query, err := r.db.Prepare("CALL SetNewPassword(?, ?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL SetNewPassword(%v, %s).", userID, password)
		return err
	}

	defer query.Close()

	res, err := query.Exec(userID, password)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to EXECUTE statement for:  CALL SetNewPassword(%v, %s).", userID, password)
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		logger.Error().Err(err).Msgf("failed get number of affected rows CALL SetNewPassword(%v, %v).", userID, password)
		return err
	}

	if ra == 0 {
		err = fmt.Errorf("rows Affected = 0 for statement: CALL SetNewPassword(%v, %v)", userID, password)
		logger.Error().Err(err).Msgf("Rows Affected = 0 for statement: CALL SetNewPassword(%v, %v)", userID, password)

		return err
	}

	return nil
}

// UpdateUser updates user profile.
func (r *Repository) UpdateUser(user *store.User) (*store.User, error) {
	query, err := r.db.Prepare("CALL UpdateUser(?, ?, ?)")
	if err != nil {
		logger.Info().Msgf("failed to prepare statement: UpdateUser(%v, %v, %v)",
			user.ID, user.Name, user.Phone)

		return nil, err
	}

	defer query.Close()

	err = query.QueryRow(user.ID, user.Name, user.Phone).
		Scan(&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Language,
			&user.Role)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New(store.UserNotFound)
	}

	if err != nil {
		logger.Error().
			Err(err).
			Msgf("There was an error executing query: CALL UpdateUser(%v, %v, %v)",
				user.ID, user.Name, user.Phone)

		return nil, err
	}

	return user, nil
}

func fetchUserRoles(rows *sql.Rows) ([]*store.Role, error) {
	userRoles := []*store.Role{}

	for rows.Next() {
		var roleName sql.NullString
		var roleID sql.NullInt64

		err := rows.Scan(&roleID, &roleName)
		if err != nil {
			logger.Error().Err(err).Msgf("failed to read role")
			return nil, err
		}

		userRole := &store.Role{
			Name:  roleName.String,
			Value: strings.ToUpper(roleName.String),
			ID:    roleID.Int64,
		}

		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}
