package http

import (
	"context"
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
)

type userKey string

const userCtx userKey = "user"

//DELETE THIS COMMENT!!!!!!!!

func (h *Handler) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			username, password, ok := r.BasicAuth()
			if !ok {
				_ = writeJSONError(w, 401, "authorization header is missing or malformed")
				return
			}

			expectedUser := os.Getenv("BASIC_AUTH_USER")
			expectedPass := os.Getenv("BASIC_AUTH_PASS")

			userMatch := subtle.ConstantTimeCompare([]byte(username), []byte(expectedUser)) == 1
			passMatch := subtle.ConstantTimeCompare([]byte(password), []byte(expectedPass)) == 1

			if !userMatch || !passMatch {
				_ = writeJSONError(w, 401, "invalid credentials")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (h *Handler) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			_ = writeJSONError(w, 401, "authorization header is missing")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			_ = writeJSONError(w, 401, "authorization header is malformed")
			return
		}
		token := parts[1]
		pasetoToken, err := h.authenticator.VerifyToken(token)
		if err != nil {
			_ = writeJSONError(w, 401, "unauthorized")
			return
		}
		ctx := r.Context()
		user, err := h.users.GetByID(ctx, pasetoToken.UserID)
		if err != nil {
			_ = writeJSONError(w, 401, "unauthorized")
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
