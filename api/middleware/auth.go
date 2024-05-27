package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/adinovcina/golang-setup/api"
	"github.com/adinovcina/golang-setup/config"
	"github.com/adinovcina/golang-setup/store"
	jwt "github.com/golang-jwt/jwt"
	"github.com/twinj/uuid"
)

type sessionFetcher interface {
	GetSession(ctx context.Context, uid uuid.UUID, sid string) (string, error)
}

func AuthorizeRequest(conf *config.Redis, inMemRepo sessionFetcher) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Request data object
			data := api.RequestData(r)
			bearerToken, err := getTokenFromHeader(r)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if bearerToken == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// 1. Validate the JWT token
			token, tokenErr := jwt.ParseWithClaims(bearerToken, &api.Claim{},
				func(token *jwt.Token) (interface{}, error) {
					// Make sure token's signature wasn't changed
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, errors.New("unexpected signing method")
					}

					return []byte(conf.SecretKey), nil
				})

			if tokenErr != nil || !token.Valid {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// 2. Get claim
			claim := token.Claims.(*api.Claim)

			// 3. Get the session object for the user
			sessionDataRaw, sessionErr := inMemRepo.GetSession(r.Context(), claim.UserID, claim.SessionID)
			if sessionErr != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// 4. Unmarshal stored object
			userData := &api.Data{}

			err = json.Unmarshal([]byte(sessionDataRaw), userData)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			data.UserID = userData.UserID
			data.Email = userData.Email
			data.Active = userData.Active
			data.Role = userData.Role
			data.UserRoleID = userData.UserRoleID
			data.SessionKey = userData.SessionKey

			// Create new context and pass new updated data
			ctx := api.NewContextWithMiddlewareData(r.Context(), data)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CheckAllowedRoles(roles ...store.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := api.RequestData(r)
			hasRole := false

			for _, role := range roles {
				if role.Name == data.Role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Create new context and pass new updated data
			ctx := api.NewContextWithMiddlewareData(r.Context(), data)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getTokenFromHeader(r *http.Request) (string, error) {
	const keyAuthorization, keyBearer, lenOfTwo = "Authorization", "Bearer", 2

	authHeader := r.Header.Get(keyAuthorization)

	authHeaderSplit := strings.Split(authHeader, " ")
	if authHeaderSplit[0] != keyBearer {
		return "", errors.New("auth err: invalid token")
	}

	if len(authHeaderSplit) < lenOfTwo {
		return "", errors.New("auth err: invalid auth header length")
	}

	return authHeaderSplit[1], nil
}
