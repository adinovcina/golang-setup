package mysqlstore

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adinovcina/golang-setup/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func (s *RepositorySuite) TestGetUserByEmail() {
	userID := uuid.NewV4()

	tests := []struct {
		name        string
		queryResult *sqlmock.Rows
		expected    *store.User
		expectErr   bool
		errorMsg    interface{}
	}{
		{
			name: "Success Case",
			queryResult: sqlmock.NewRows([]string{
				"ID", "Name", "Email", "Password",
				"Active", "FailedLoginCount", "LoginBlockedUntil",
			}).
				AddRow(
					userID, "test user", "test@gmail.com", "$2a$10$HnQIEV5YpB8BxXjr6p5UuuVo901a/W/fHo3GDHbslZw1RZvYsPtWG",
					true, 0, nil,
				),
			expected: &store.User{
				ID:                userID,
				Name:              "test user",
				Email:             "test@gmail.com",
				Password:          "$2a$10$HnQIEV5YpB8BxXjr6p5UuuVo901a/W/fHo3GDHbslZw1RZvYsPtWG",
				Active:            true,
				FailedLoginCount:  0,
				LoginBlockedUntil: nil,
			},
			expectErr: false,
			errorMsg:  nil,
		},
		{
			name:        "Error Case - User not found",
			queryResult: sqlmock.NewRows([]string{}),
			expected: &store.User{
				Email: "invalidEmail@gmail.com",
			},
			expectErr: true,
			errorMsg:  store.UserNotFound,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.mock.ExpectPrepare("^CALL GetUserByEmail\\(\\?\\)$").
				ExpectQuery().
				WithArgs(tt.expected.Email).
				WillReturnRows(tt.queryResult)

			user, err := s.repo.GetUserByEmail(tt.expected.Email)

			if tt.expectErr {
				require.Error(t, err)
				require.Equal(t, tt.errorMsg, store.UserNotFound)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, user)
				require.Equal(t, tt.errorMsg, nil)
			}

			err = s.mock.ExpectationsWereMet()
			s.Require().NoError(err)
		})
	}
}

func (s *RepositorySuite) TestGetUserByID() {
	currentTime := time.Now()
	userID := uuid.NewV4()

	tests := []struct {
		name           string
		queryResult    *sqlmock.Rows
		expected       *store.User
		expectErr      bool
		errorMsg       interface{}
		logExpectation func(mockLogger *zerolog.Event)
	}{
		{
			name: "Success Case",
			queryResult: sqlmock.NewRows([]string{
				"ID", "Name", "Email", "Phone",
				"Language", "Active", "Role", "CreatedAt",
			}).
				AddRow(
					userID, "test user", "test@gmail.com", "+12312313",
					"en", true, "user", currentTime,
				),
			expected: &store.User{
				ID:        userID,
				Name:      "test user",
				Email:     "test@gmail.com",
				Phone:     "+12312313",
				Language:  "en",
				Role:      "user",
				CreatedAt: currentTime,
				Active:    true,
			},
			expectErr:      false,
			errorMsg:       nil,
			logExpectation: nil,
		},
		{
			name:        "Error Case - User not found",
			queryResult: sqlmock.NewRows([]string{}),
			expected: &store.User{
				ID: uuid.Nil.UUID(),
			},
			expectErr: true,
			errorMsg:  store.UserNotFound,
			logExpectation: func(mockLogger *zerolog.Event) {
				mockLogger.Msgf("failed to prepare statement: CALL GetUserByID(%s).", uuid.Nil.UUID())
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.mock.ExpectPrepare("^CALL GetUserByID\\(\\?\\)$").
				ExpectQuery().
				WithArgs(tt.expected.ID).
				WillReturnRows(tt.queryResult)

			user, err := s.repo.GetUserByID(tt.expected.ID)

			if tt.expectErr {
				require.Error(t, err)
				require.Equal(t, tt.errorMsg, store.UserNotFound)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, user)
				require.Equal(t, tt.errorMsg, nil)
			}

			err = s.mock.ExpectationsWereMet()
			s.Require().NoError(err)
		})
	}
}
