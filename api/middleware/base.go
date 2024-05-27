package middleware

import (
	"net/http"
	"strings"

	"github.com/adinovcina/golang-setup/api"
	"github.com/adinovcina/golang-setup/tools/utils"
)

func InitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First check if the data has already been initialized
		rx := api.MiddlewareDataFromContext(r.Context())
		if rx != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Initialize with default values. It is to be updated later on inside other middlewares and handlers.
		// Make sure that this middleware is executed first
		requestID := r.Header.Get("X-Request-Id")

		// There is no requestID, then add one manually double uuid as requestID
		if requestID == "" {
			requestID = strings.ReplaceAll(utils.GenerateUniqueID()+utils.GenerateUniqueID(), "-", "")
		}

		data := &api.Data{
			RequestID: requestID,
		}

		// Add cors
		w.Header().Set("Content-Type", "application/json charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Request-Method", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Expose-Headers",
			`App-Token, Status-Code, X-Request-Id, X-Content-Length, Content-Length`)
		w.Header().Set(`Access-Control-Allow-Headers`, `Origin, X-Requested-With, Content-Type, X-Content-Length, 
			Content-Length, Accept-Encoding, Accept, Access-Control-Allow-Origin, Authorization, App-Token, Status-Code, 
			Access-Control-Allow-Credentials,X-Request-Id,Access-Control-Request-Method`)

		if r.Method == http.MethodOptions {
			return
		}

		ctx := api.NewContextWithMiddlewareData(r.Context(), data)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
