package http

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"os"
	"strings"
)

type userKey string

const userCtx userKey = "user"

func (h *Handler) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSONError(w, 401, "authorization header is missing")
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				writeJSONError(w, 401, "authorization header is malformed")
				return
			}
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				writeJSONError(w, 401, err.Error())
				return
			}
			username := os.Getenv("BASIC_AUTH_USER")
			password := os.Getenv("BASIC_AUTH_PASS")
			creds := strings.SplitN(string(decoded), ":", 2)
			userMatch := subtle.ConstantTimeCompare([]byte(creds[0]), []byte(username)) == 1
			passMatch := subtle.ConstantTimeCompare([]byte(creds[1]), []byte(password)) == 1

			if !userMatch || !passMatch {
				writeJSONError(w, 401, "invalid credentials")
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
			writeJSONError(w, 401, "authorization header is missing")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeJSONError(w, 401, "authorization header is malformed")
			return
		}
		token := parts[1]
		pasetoToken, err := h.authenticator.VerifyToken(token)
		if err != nil {
			writeJSONError(w, 401, "unauthorized")
			return
		}
		ctx := r.Context()
		user, err := h.users.GetByID(ctx, pasetoToken.UserID)
		if err != nil {
			writeJSONError(w, 401, err.Error())
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
