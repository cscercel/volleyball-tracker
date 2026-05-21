package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/cscercel/volleyball-tracker/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MatchHandler struct {
	service	*service.MatchService
}

func NewMatchHandler(service *service.MatchService) *MatchHandler {
	return &MatchHandler{service: service}
}

func (h *MatchHandler) RegisterRoutes(r chi.Router) {
	r.Route("/matches", func(r chi.Router) {
		r.Post("/", h.handleCreateMatch)
		r.Get("/{id}", h.handleGetMatch)
		r.Get("/registered", h.handleGetRegisteredMatches)
		r.Get("/drafts", h.handleGetDrafts)
		r.Get("/", h.handleGetSeasonMatches)
		r.Put("/{id}", h.handleRegisterMatch)
		r.Delete("/{id}", h.handleDeleteDraft)
	})
}


// @Summary      Get Match
// @Tags         matches
// @Produce      json
// @Param        id   path      string  true  "Match ID"
// @Success      200  {array}   service.MatchWithPlayers
// @Failure      400  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /api/v1/matches/{id} [get]
func (h *MatchHandler) handleGetMatch(w http.ResponseWriter, r *http.Request) {
	matchID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid match id", err)
		return
	}

	match, err := h.service.GetMatch(r.Context(), matchID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "match not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, match)
}

// @Summary      Get Registered Matches
// @Tags         matches
// @Produce      json
// @Success      200  {array}   []db.Match
// @Failure      400  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /api/v1/matches/registered [get]
func (h *MatchHandler) handleGetRegisteredMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := h.service.GetRegisteredMatches(r.Context())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "matches not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, matches)
}

// @Summary      Get Drafts
// @Tags         matches
// @Produce      json
// @Success      200  {array}   []db.Match
// @Failure      400  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /api/v1/matches/drafts [get]
func (h *MatchHandler) handleGetDrafts(w http.ResponseWriter, r *http.Request) {
	drafts, err := h.service.GetRegisteredMatches(r.Context())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "drafts not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, drafts)
}

// @Summary      Get Season Matches
// @Tags         matches
// @Produce      json
// @Param        match_type  path     string  true  "Match Type"
// @Param        season  	path      string  true  "Season"
// @Success      200  		{array}   []db.Match
// @Failure      400  		{object}  object{error=string}
// @Failure      404  		{object}  object{error=string}
// @Router       /api/v1/matches [get]
func (h *MatchHandler) handleGetSeasonMatches(w http.ResponseWriter, r *http.Request) {
	// Query params
	matchType := r.URL.Query().Get("match_type")
	if matchType == "" {
		respondWithError(w, http.StatusBadRequest, "match_type is required", errors.New(""))
		return
	}

	seasonStr := r.URL.Query().Get("season")
	season, err := strconv.Atoi(seasonStr)
	if err != nil || season < 1 {
		respondWithError(w, http.StatusBadRequest, "season must be a positive number", err)
		return
	}

	matches, err := h.service.GetSeasonMatches(r.Context(), matchType, season)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "drafts not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, matches)
}

// @Summary      Create Match
// @Tags         matches
// @Produce      json
// @Param        body body      object{match_type=string, blue_players={array}, red_players={array}} true "Match Body"
// @Success      201  {array}   service.MatchWithPlayers
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/matches [post]
func (h *MatchHandler) handleCreateMatch(w http.ResponseWriter, r *http.Request) {
	var body struct {
		MatchType	string 	 	`json:"match_type"`
		BluePlayers	[]uuid.UUID `json:"blue_players"`
		RedPlayers	[]uuid.UUID `json:"red_players"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Get

	match, err := h.service.CreateMatch(r.Context(), body.MatchType, body.BluePlayers, body.RedPlayers)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create match", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, match)
}

// @Summary      Register Match
// @Tags         matches
// @Produce      json
// @Param        id   path      string                				  true  "Match UUID" format(uuid)
// @Param        body body      object{blue_score=int, red_score=int} true "Match Body"
// @Success      200  {array}   db.Match
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/matches/{id} [put]
func (h *MatchHandler) handleRegisterMatch(w http.ResponseWriter, r *http.Request) {
	matchID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid match id", err)
		return
	}

	var body struct {
		MatchID		uuid.UUID 	`json:"match_id"`
		BlueScore 	int			`json:"blue_score"`
		RedScore	int 		`json:"red_score"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	body.MatchID = matchID

	match, err := h.service.RegisterMatch(r.Context(), body.MatchID, body.BlueScore, body.RedScore)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not register match", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, match)
}

// @Summary      Delete Match
// @Tags         matches
// @Produce      json
// @Param        id   path      string                true  "Match UUID" format(uuid)
// @Success      204  
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/matches/{id} [delete]
func (h *MatchHandler) handleDeleteDraft(w http.ResponseWriter, r *http.Request) {
	matchID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid match id", err)
		return
	}

	if err := h.service.DeleteDraft(r.Context(), matchID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete draft", err)
		return
	}
	
	respondWithJSON(w, http.StatusNoContent, "")
}
