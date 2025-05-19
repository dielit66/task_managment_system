package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dielit66/task-management-system/internal/logger"
)

type JWTClaims struct {
	UserID int `json:"user_id"`
}

func JwtPayloadMiddleware(l logger.ILogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				l.Error("Authorization header is missing")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				l.Error("Invalid authorization header format")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := tokenParts[1]
			claims, err := parsePayloadFromJWT(token)
			if err != nil {
				l.Error("Failed to parse JWT", "token", token, "error", err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if claims.UserID == 0 {
				l.Error("Invalid user_id in JWT claims", "user_id", claims.UserID)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			l.Debug("JWT parsed successfully", "user_id", claims.UserID)
			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parsePayloadFromJWT(t string) (*JWTClaims, error) {
	var claims JWTClaims

	tokenParts := strings.Split(t, ".")
	if len(tokenParts) != 3 {
		return nil, fmt.Errorf("invalid JWT format: expected 3 parts")
	}

	payload := tokenParts[1]
	dStr, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("base64url decode error: %v", err)
	}

	err = json.Unmarshal(dStr, &claims)
	if err != nil {
		return nil, fmt.Errorf("JSON unmarshal error: %v", err)
	}

	return &claims, nil
}
