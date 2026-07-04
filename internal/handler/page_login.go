package handler

import (
	"net/http"
	"time"

	"github.com/cscercel/volleyball-tracker/web/templates/pages"
)

func (h *PageHandler) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	pages.Login("", isAuthenticated(r, h.jwtSecret)).Render(r.Context(), w)
}

func (h *PageHandler) handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		pages.LoginError("Could not read form data.").Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.userService.GetUser(r.Context(), email, password)
	if err != nil {
		pages.LoginError("Invalid email or password").Render(r.Context(), w)
		return
	}

	setAuthCookie(w, user.Token, time.Hour, h.secureCookies)

	w.Header().Set("HX-Redirect", "/")
}

func (h *PageHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	clearAuthCookie(w, h.secureCookies)
	w.Header().Set("HX-Redirect", "/login")
}
