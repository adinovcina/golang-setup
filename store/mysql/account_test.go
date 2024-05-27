package mysqlstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adinovcina/golang-setup/store"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func (s *RepositorySuite) TestGetUserRoles() {
	userID := uuid.NewV4()

	tests := []struct {
		name        string
		queryResult *sqlmock.Rows
		queryParam  uuid.UUID
		expected    []*store.Role
		errorResp   interface{}
	}{
		{
			name: "Success Case",
			queryResult: sqlmock.NewRows([]string{
				"ID", "Name",
			}).AddRow(
				1, "Admin",
			),
			queryParam: userID,
			expected: []*store.Role{
				{
					ID:    1,
					Name:  "Admin",
					Value: "ADMIN",
				},
			},
			errorResp: nil,
		},
		{
			name:        "Error Case - The user does not have roles assigned",
			queryResult: sqlmock.NewRows([]string{}),
			expected:    []*store.Role{},
			errorResp:   nil,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.mock.ExpectPrepare("^CALL GetRoles\\(\\?\\)$").
				ExpectQuery().
				WithArgs(tt.queryParam).
				WillReturnRows(tt.queryResult)

			roles, err := s.repo.GetUserRoles(tt.queryParam)

			require.NoError(t, err)
			require.Equal(t, tt.expected, roles)
			require.Equal(t, tt.errorResp, nil)

			err = s.mock.ExpectationsWereMet()
			s.Require().NoError(err)
		})
	}
}
