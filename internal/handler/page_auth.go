package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/cscercel/volleyball-tracker/internal/auth"
)

const authCookieName = "auth_token"

func setAuthCooke(w http.ResponseWriter, token string, expiresIn time.Duration, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name: authCookieName,
		Value: token,
		Path: "/",
		HttpOnly: true,
		Secure: secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge: int(expiresIn.Seconds()),
	})
}

func clearAuthCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name: authCookieName,
		Value: "",
		Path: "/",
		HttpOnly: true,
		Secure: secure,
		SameSite: http.SameSiteDefaultMode,
		MaxAge: -1, // deletes the cookie instant speed
	})
}

func PageAuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(authCookieName)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			userID, err := auth.ValidateJWT(cookie.Value, jwtSecret)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), contextKey("userID"), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
