package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cscercel/volleyball-tracker/internal/auth"
	_ "github.com/cscercel/volleyball-tracker/internal/db" // ONLY required for Swagger to pick up db interfaces
	"github.com/cscercel/volleyball-tracker/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	service		*service.UserService
	jwtSecret	string
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.handleCreateUser)
		r.Get("/{id}", h.handleLogin)
		r.Put("/{id}/change-email", h.handleUpdateUserEmail)
		r.Put("/{id}/password-reset", h.handleUpdateUserPassword)
	})
}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email				string	`json:"email"`
		Password			string	`json:"password"`
		RegistrationCode	string	`json:"registration_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	user, err := h.service.CreateUser(r.Context(), body.Email, body.Password, body.RegistrationCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email		string	`json:"email"`
		Password	string	`json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	user, err := h.service.GetUser(r.Context(), body.Email, body.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) handleUpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	var body struct {
		NewEmail	string	`json:"new_email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Check if user is authenticated
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token in header", err)
		return
	}

	id, err := auth.ValidateJWT(token, h.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate token", err)
		return
	}

	user, err := h.service.UpdateUserEmail(r.Context(), id, body.NewEmail)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not change email", err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) handleUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		NewPassword	string	`json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Check if user is authenticated
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token in header", err)
		return
	}

	id, err := auth.ValidateJWT(token, h.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate token", err)
		return
	}

	user, err := h.service.UpdateUserPassword(r.Context(), id, body.NewPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not change password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
