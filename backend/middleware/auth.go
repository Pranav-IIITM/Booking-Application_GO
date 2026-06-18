package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

type contextKey string

const (
	firebaseUIDKey contextKey = "firebaseUID"
	emailKey       contextKey = "email"
	nameKey        contextKey = "name"
)

func FirebaseAuth(authClient *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.Fields(authHeader)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeError(w, http.StatusUnauthorized, "authorization header must be Bearer token")
				return
			}

			token, err := authClient.VerifyIDToken(r.Context(), parts[1])
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), firebaseUIDKey, token.UID)
			if email, ok := token.Claims["email"].(string); ok {
				ctx = context.WithValue(ctx, emailKey, email)
			}
			if name, ok := token.Claims["name"].(string); ok {
				ctx = context.WithValue(ctx, nameKey, name)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FirebaseUID(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(firebaseUIDKey).(string)
	return uid, ok && uid != ""
}

func Email(ctx context.Context) string {
	email, _ := ctx.Value(emailKey).(string)
	return email
}

func Name(ctx context.Context) string {
	name, _ := ctx.Value(nameKey).(string)
	return name
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
