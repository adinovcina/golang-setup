package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adinovcina/golang-setup/api"
	"github.com/adinovcina/golang-setup/config"
	"github.com/adinovcina/golang-setup/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func TestCheckAllowedRoles(t *testing.T) {
	// Define some sample roles
	adminRole := store.Role{Name: "Admin"}
	userRole := store.Role{Name: "User"}

	// Create a test handler that always returns OK
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Define test cases
	tests := []struct {
		name           string
		roles          []store.Role
		requestRole    string
		expectedStatus int
	}{
		{
			name:           "Allowed Role",
			roles:          []store.Role{adminRole, userRole},
			requestRole:    "Admin",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Denied Role",
			roles:          []store.Role{adminRole},
			requestRole:    "User",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "No Roles Provided",
			roles:          nil,
			requestRole:    "User",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new request with the specified role
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req = req.WithContext(api.NewContextWithMiddlewareData(req.Context(), &api.Data{Role: tc.requestRole}))

			// Create a new recorder to record the response
			rr := httptest.NewRecorder()

			// Create a middleware with the specified roles
			middleware := CheckAllowedRoles(tc.roles...)

			// Execute the middleware chain with the test handler
			middleware(testHandler).ServeHTTP(rr, req)

			// Check if the status code matches the expected status code
			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}

type mockSessionFetcher struct{}

func (m *mockSessionFetcher) GetSession(ctx context.Context, uid uuid.UUID, sid string) (string, error) {
	// Mock implementation for session fetching
	return `{"userID": "0a15f901-55a7-4dac-b1ae-c602fb775bd1", "email": "admin@gmail.com", "active": true, "role": "Admin", "userRoleID": 1, "sessionKey": "sessionKey"}`, nil
}

func TestAuthorizeRequest(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		token    string
		expected int
	}{
		{
			name:     "ValidToken",
			token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI5MTU3MDY3NjYsInNlc3Npb25JRCI6Ii1OeG13SlVzd29JRDdMV0NxTE5DIiwidXNlcklEIjoiMGExNWY5MDEtNTVhNy00ZGFjLWIxYWUtYzYwMmZiNzc1YmQxIn0.wPJ2B8V-1-Fym65xy0PCWLHRDsGncVZ1XRNg8ljWxu0",
			expected: http.StatusOK,
		},
		{
			name:     "InvalidToken",
			token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTU3MDY3NjYsInNlc3Npb25JRCI6Ii1OeG13SlVzd29JRDdMV0NxTE5DIiwidXNlcklEIjoiMGExNWY5MDEtNTVhNy00ZGFjLWIxYWUtYzYwMmZiNzc1YmQxIn0.hxHKaTSzzdfQpz-ElJiQwcryyJ60jQQ5Wr1t_NBBXU4",
			expected: http.StatusUnauthorized,
		},
		{
			name:     "MissingToken",
			token:    "",
			expected: http.StatusUnauthorized,
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare a mock Redis configuration
			mockConf := &config.Redis{
				SecretKey: "test",
			}

			// Create a mock session fetcher
			mockSession := &mockSessionFetcher{}

			// Create a mock HTTP handler
			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Retrieve data from context
				data := api.MiddlewareDataFromContext(r.Context())

				// Check if data is nil
				if data == nil {
					t.Error("Middleware data is nil")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
			})

			// Create a request with the token
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tc.token)

			userID, _ := uuid.Parse("0a15f901-55a7-4dac-b1ae-c602fb775bd1")

			// Create context with middleware data
			ctx := api.NewContextWithMiddlewareData(req.Context(), &api.Data{
				UserID:     *userID,
				Email:      "admin@gmail.com",
				Active:     true,
				Role:       "Admin",
				UserRoleID: 1,
				SessionKey: "sessionKey",
			})

			// Add context to the request
			req = req.WithContext(ctx)

			// Create a recorder to record the response
			rr := httptest.NewRecorder()

			// Call the middleware with the mock session fetcher
			AuthorizeRequest(mockConf, mockSession)(mockHandler).ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tc.expected, rr.Code)
		})
	}
}
