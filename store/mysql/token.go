package mysqlstore

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/adinovcina/golang-setup/store"
	"github.com/adinovcina/golang-setup/tools/logger"
	"github.com/twinj/uuid"
)

// GetPasswordTokenByToken retrieves password token by token.
func (r *Repository) GetPasswordTokenByToken(token string) (*store.PasswordToken, error) {
	query, err := r.db.Prepare("CALL GetPasswordTokenByToken(?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL GetPasswordTokenByToken(%s)", token)
		return nil, err
	}

	defer query.Close()

	passwordToken := new(store.PasswordToken)

	err = query.QueryRow(token).Scan(&passwordToken.UserID, &passwordToken.Token, &passwordToken.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("TOKEN_NOT_FOUND")
	}

	if err != nil {
		logger.Error().Err(err).Msgf("There was an error executing query: CALL GetPasswordTokenByToken(%s)", token)
		return nil, err
	}

	return passwordToken, nil
}

// AddPasswordResetToken add new reset password for the user.
func (r *Repository) AddPasswordResetToken(userID uuid.UUID, token string, expiresAt int64) (*store.PasswordToken, error) {
	query, err := r.db.Prepare("CALL AddPasswordToken(?, ?, ?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL AddPasswordToken(%v, %v, %v)",
			userID, token, expiresAt)
		return nil, err
	}

	defer query.Close()

	passwordToken := new(store.PasswordToken)

	err = query.QueryRow(userID, token, expiresAt).
		Scan(&passwordToken.UserID,
			&passwordToken.Token,
			&passwordToken.ExpiresAt)
	if err != nil {
		logger.Error().Err(err).Msgf("There was an error executing query: CALL AddPasswordToken(%v, %v, %v)",
			userID, token, expiresAt)
		return nil, err
	}

	return passwordToken, nil
}

// GetTokenByTokenAndType will retrieve the token by token value and type.
func (r *Repository) GetTokenByTokenAndType(token, tokenType string) (*store.LoginToken, error) {
	query, err := r.db.Prepare("CALL GetTokenByTokenAndType(?,?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL GetTokenByTokenAndType(%v,%v).",
			token, tokenType)
		return nil, err
	}

	defer query.Close()

	model := new(store.LoginToken)

	err = query.QueryRow(token, tokenType).
		Scan(&model.Expired,
			&model.ID,
			&model.UserID,
			&model.Token,
			&model.TokenType)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New(store.TokenNotFound)
	}

	if err != nil {
		logger.Error().Err(err).Msgf("There was an error executing query:CALL GetTokenByTokenAndType(%v,%v).",
			token, tokenType)
		return nil, err
	}

	return model, nil
}

// DeleteTokenByID deletes token by id.
func (r *Repository) DeleteTokenByID(id int64) error {
	query, err := r.db.Prepare("CALL DeleteTokenByID(?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL DeleteTokenByID(%v)", id)
		return err
	}

	defer query.Close()

	res, err := query.Exec(id)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to execute statement: CALL DeleteTokenByID(%v)", id)
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		logger.Error().Err(err).Msgf("failed after statement is executed: CALL DeleteTokenByID(%v)", id)
		return err
	}

	if ra == 0 {
		err = fmt.Errorf("rows Affected = 0 for statement: CALL  DeleteTokenByID(%v)", id)
		logger.Error().Err(err).Msgf("No rows were deleted: CALL DeleteTokenByID(%v)", id)

		return err
	}

	return nil
}
