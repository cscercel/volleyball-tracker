package handler

import (
	"encoding/json"
	"net/http"

	_ "github.com/cscercel/volleyball-tracker/internal/db" // ONLY required for Swagger to pick up db interfaces
	"github.com/cscercel/volleyball-tracker/internal/service"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		// Public routes
		r.Post("/", h.handleCreateUser)
		r.Post("/login", h.handleLogin)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Put("/{id}/change-email", h.handleUpdateUserEmail)
			r.Put("/{id}/password-reset", h.handleUpdateUserPassword)
		})
	})
}

// @Summary      Create User
// @Tags         users
// @Produce      json
// @Param        body body      object{email=string,password=string,registration_code=string} true "User Body"
// @Success      201  {array}   db.CreateUserRow
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/users [post]
func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		RegistrationCode string `json:"registration_code"`
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

// @Summary      Login User
// @Tags         users
// @Produce      json
// @Param        body body      object{email=string,password=string} true "User Body"
// @Success      200  {array}   service.UserLogin
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/users/login [post]
func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

// @Summary      Update User Email
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        body body      object{email=string} true "User Body"
// @Success      200  {array}   db.UpdateUserEmailRow
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/users/change-email [put]
func (h *UserHandler) handleUpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	var body struct {
		NewEmail string `json:"new_email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	id, err := GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	user, err := h.service.UpdateUserEmail(r.Context(), id, body.NewEmail)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not change email", err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// @Summary      Update User Password
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        body body      object{password=string} true "User Body"
// @Success      200  {array}   db.UpdateUserPasswordRow
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/users/password-reset [put]
func (h *UserHandler) handleUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	id, err := GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	user, err := h.service.UpdateUserPassword(r.Context(), id, body.NewPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not change password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
