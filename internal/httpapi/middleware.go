package httpapi

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/oopchoop/litebank/internal/token"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(jwtManager *token.JWTManager) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeError(w, http.StatusUnauthorized, "invalid authorization format")
				return
			}

			tokenString := parts[1]

			userID, err := jwtManager.Parse(tokenString)
			if err != nil {
				writeError(w, http.StatusUnauthorized, publicMessage(err))
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)

	return id, ok
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic error: %v\n%s", err, debug.Stack())

				writeError(w, http.StatusInternalServerError, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
