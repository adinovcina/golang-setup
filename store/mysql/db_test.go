package mysqlstore

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

// RepositorySuite is a test suite for db package.
type RepositorySuite struct {
	suite.Suite
	repo *Repository
	mock sqlmock.Sqlmock
}

// TestProviderSuite is a suite for area group unit tests.
func TestProviderSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositorySuite))
}

// SetupSuite configures suite for unit testing.
func (s *RepositorySuite) SetupSuite() {
	// initializes a mocked DB and a mock object
	testDB, mock, err := sqlmock.New()
	s.mock = mock
	repo := New(testDB, nil)
	s.repo = repo
	s.Require().NoError(err)
	s.Require().NotNil(repo)
}

// TearDownSuite ensures that resources are cleaned up after all tests in suite are run.
func (s *RepositorySuite) TearDownSuite() {
	s.repo.db.Close()
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface.
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
