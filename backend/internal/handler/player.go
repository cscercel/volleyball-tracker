package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/cscercel/volleyball-tracker/internal/service"
)

type PlayerHandler struct {
	playerService	*service.PlayerService
}

func (h *PlayerHandler) GetPlayerCareer(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	player, err := h.playerService.GetPlayerCareer(r.Context(), id)
}
