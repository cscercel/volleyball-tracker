package handler

import (
	"net/http"
)

type PlayerHandler struct {
	playerService	*service.PlayerService
}

func (h *PlayerHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	
}
