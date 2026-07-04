package handler

import (
	"net/http"

	"github.com/cscercel/volleyball-tracker/web/templates/pages"
)

func (h *PageHandler) handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	pages.Register(isAuthenticated(r, h.jwtSecret)).Render(r.Context(), w)
}

func (h *PageHandler) handleRegisterSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		pages.RegisterStatus("Could not read form data.", false).Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	registrationCode := r.FormValue("registration_code")

	if _, err := h.userService.CreateUser(r.Context(), email, password, registrationCode); err != nil {
		pages.RegisterStatus("Could not create account.", false).Render(r.Context(), w)
		return
	}

	pages.RegisterStatus("Account created", true).Render(r.Context(), w)
}
