package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cscercel/volleyball-tracker/internal/auth"
	"github.com/google/uuid"
)

type contextKey string

func AuthenticateMiddleware(tokenSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.GetBearerToken(r.Header)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "failed to retrieve Bearer token", err)
				return
			}

			userID, err := auth.ValidateJWT(token, tokenSecret)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "failed to authenticate user", err)
				return
			}

			ctx := context.WithValue(r.Context(), contextKey("userID"), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(contextKey("userID")).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}

	return userID, nil
}
