package handlers

import (
	"io"
	"net/http"

	"github.com/adinovcina/golang-setup/api/account"
	m "github.com/adinovcina/golang-setup/api/middleware"
	"github.com/adinovcina/golang-setup/config"
	"github.com/adinovcina/golang-setup/services"
	"github.com/adinovcina/golang-setup/store"
	"github.com/adinovcina/golang-setup/tools/logger"
	s "github.com/adinovcina/golang-setup/tools/network/http"
	"github.com/go-chi/chi/v5"
)

// healthCheck method is used to check if server is live.
func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	if _, err := io.WriteString(w, "OK"); err != nil {
		logger.Error().Err(err).Msg("write string failed")
	}
}

// Attach is used to setup.
func Attach(server *s.Server,
	repo store.Repository,
	inMemRepo store.InMemRepository,
	conf *config.Config,
	appServices *services.AppServices,
) *s.Server {
	// Apply default unprotected middlewares to root api group
	publicGroup := server.Get().Route("/", func(r chi.Router) {
		r.Use(m.InitMiddleware)
		r.Use(m.Logger)
	})

	// Apply protected middleware to group
	_ = publicGroup.Route("/", func(r chi.Router) {
		r.Use(m.AuthorizeRequest(&conf.Redis, inMemRepo))
	})

	// Health check route.
	publicGroup.Get("/health", healthCheck)

	// Attach Account Routes.
	account.AttachAccountRoutes(publicGroup,
		conf,
		repo,
		inMemRepo,
		appServices.GetMailjetClient())

	return server
}
