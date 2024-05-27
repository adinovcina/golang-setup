package middleware

import (
	"net/http"
	"time"

	"github.com/adinovcina/golang-setup/api"
	"github.com/adinovcina/golang-setup/tools/logger"
	"github.com/adinovcina/golang-setup/tools/utils"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const (
			statusCode = "status"
			httpMethod = "method"
			path       = "path"
			host       = "host"
			reqQuery   = "request_query"
			latency    = "latency"
			userAgent  = "user_agent"
			httpReq    = "http_request"
			requestID  = "request_id"
		)

		if r.Method == http.MethodOptions {
			return
		}

		if r.RequestURI == "/health" || r.RequestURI == "/favicon.ico" || r.RequestURI == "/" {
			return
		}

		start := time.Now()

		// Do not delay the request, let it pass.
		// Deal with rest of the data in an immutable fashion!
		sw := new(statusWriter)
		sw.ResponseWriter = w

		next.ServeHTTP(sw, r)

		if shouldOmitLog(r.RequestURI) {
			return
		}

		requestData := api.RequestData(r)

		reqLogger := logger.With().
			Int(statusCode, sw.status).
			Str(httpMethod, r.Method).
			Str(path, r.RequestURI).
			Str(host, r.Host). // Note that this will rarely produce the correct IP address.
			Str(latency, time.Since(start).String()).
			Str(userAgent, r.Header.Get("User-Agent")).
			Str(requestID, requestData.RequestID).
			Logger()

		switch {
		case sw.status >= http.StatusBadRequest && sw.status < http.StatusInternalServerError:
			reqLogger.Warn().Msg(httpReq)
		case sw.status >= http.StatusInternalServerError:
			reqLogger.Error().Msg(httpReq)
		default:
			reqLogger.Info().Msg(httpReq)
		}
	})
}

func shouldOmitLog(requestURI string) bool {
	requestURIsToOmit := []string{"/", "/health", "/favicon.ico", "/listings/seatgeek/purchase-eligible"}

	return utils.Contains(requestURIsToOmit, requestURI)
}
