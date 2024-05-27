package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adinovcina/golang-setup/config"
	"github.com/adinovcina/golang-setup/tools/network"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	conf := &config.Config{
		Service: config.Service{
			Port: "8080",
		},
		Timeouts: config.Timeouts{
			ReadDuration:  30 * time.Second,
			WriteDuration: 30 * time.Second,
		},
	}
	server := NewServer(conf)

	assert.NotNil(t, server)
	assert.Equal(t, ":"+conf.Service.Port, server.Configuration.HTTPPort)
	assert.Equal(t, conf.Timeouts.ReadDuration, server.Configuration.HTTPReadTimeout)
	assert.Equal(t, conf.Timeouts.WriteDuration, server.Configuration.HTTPWriteTimeout)
	assert.NotNil(t, server.router)
	assert.NotNil(t, server.server)
}

func TestServerServe(t *testing.T) {
	// Create a new server instance
	server := &Server{
		router: chi.NewRouter(),
		server: &http.Server{
			Addr:              ":5050",
			Handler:           chi.NewRouter(),
			ReadHeaderTimeout: 5 * time.Second,
		},
		Configuration: network.Config{
			HTTPPort: ":5050",
		},
	}

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Create a channel to signal that the request has been handled
	done := make(chan bool)

	// Serve the request in a goroutine
	go func() {
		server.Get().ServeHTTP(rr, req)
		done <- true
	}()

	// Wait for the server to start serving requests
	time.Sleep(100 * time.Millisecond)

	// Close the server gracefully
	err := server.Close()
	require.NoError(t, err)

	// Wait for the request to be handled
	<-done

	// Verify the response status code
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
