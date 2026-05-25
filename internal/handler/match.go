package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/cscercel/volleyball-tracker/internal/db" // ONLY required for Swagger to pick up db interfaces
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
		r.Get("/", h.handleListMatchesBySeason)
		r.Get("/uncompleted", h.handleListUncompletedMatches)
		r.Delete("/{id}", h.handleDeleteUncompletedMatch)
		r.Get("/{id}/roster", h.handleGetMatchPlayers)
	})
}

// @Summary      Create Match
// @Tags         matches
// @Produce      json
// @Param        body body      object{match_type=string, blue_team=[]string, red_team=[]string} true "Match Body"
// @Success      201  {array}   db.Match
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/matches [post]
func (h *MatchHandler) handleCreateMatch(w http.ResponseWriter, r *http.Request) {
	var body struct {
		MatchType 	string 		`json:"match_type"`
		BlueTeam	[]string	`json:"blue_team"`
		RedTeam		[]string	`json:"red_team"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	match, err := h.service.CreateMatch(r.Context(), body.MatchType, body.BlueTeam, body.RedTeam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create match", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, match)
}

// @Summary      Get Match
// @Tags         matches
// @Produce      json
// @Param        id   path      string                true  "Match UUID" format(uuid)
// @Success      200  {array}   db.Match
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/matches/{id} [get]
func (h *MatchHandler) handleGetMatch(w http.ResponseWriter, r *http.Request) {
	matchID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid match id", err)
		return
	}

	match, err := h.service.GetMatch(r.Context(), matchID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get match", err)
		return
	}

	respondWithJSON(w, http.StatusOK, match)
}

// @Summary      List Matches By Season
// @Tags         matches
// @Produce      json
// @Param        match_type query  		string  true  "Match Type"
// @Param        season   	query      	integer true  "Season"
// @Success      200  		{array}   []db.Match
// @Failure      400  		{object}  object{error=string}
// @Failure      500  		{object}  object{error=string}
// @Router       /api/v1/matches [get]
func (h *MatchHandler) handleListMatchesBySeason(w http.ResponseWriter, r *http.Request) {
	// Query params
	matchType := r.URL.Query().Get("match_type")
	if matchType != "indoor" && matchType != "beach" {
		respondWithError(w, http.StatusBadRequest, "invalid match type", 
			fmt.Errorf("expected: `indoor` or `beach`, got: `%s`", matchType),
		)
		return
	}

	seasonStr := r.URL.Query().Get("season")
	season, err := strconv.Atoi(seasonStr)
	if err != nil || season < 2023 {
		respondWithError(w, http.StatusBadRequest, "no season before 2023", err)
		return
	}

	matches, err := h.service.ListMatchesBySeason(r.Context(), matchType, int32(season))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get matches", err)
		return
	}

	respondWithJSON(w, http.StatusOK, matches)
}

// @Summary      List Uncompleted Matches
// @Tags         matches
// @Produce      json
// @Success      200  		{array}   []db.Match
// @Failure      400  		{object}  object{error=string}
// @Failure      500  		{object}  object{error=string}
// @Router       /api/v1/matches/uncompleted [get]
func (h *MatchHandler) handleListUncompletedMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := h.service.ListUncompletedMatches(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not retrieve uncompleted matches", err)
		return
	}

	respondWithJSON(w, http.StatusOK, matches)
}

// @Summary      Delete Uncompleted Match
// @Tags         matches
// @Produce      json
// @Param        id   path      string                true  "Match UUID" format(uuid)
// @Success      204  
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/matches/{id} [delete]
func (h *MatchHandler) handleDeleteUncompletedMatch(w http.ResponseWriter, r *http.Request) {
	matchID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid match id", err)
		return
	}

	if err := h.service.DeleteUncompletedMatch(r.Context(), matchID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete match", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}

// @Summary      Get Match Players
// @Tags         matches
// @Produce      json
// @Param        id   path      string                true  "Match UUID" format(uuid)
// @Success      200  {array}   []db.GetMatchPlayersRow
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/matches/{id}/roster [get]
func (h *MatchHandler) handleGetMatchPlayers(w http.ResponseWriter, r *http.Request) {
	matchID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid match id", err)
		return
	}

	match_players, err := h.service.GetMatchPlayers(r.Context(), matchID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get match players", err)
		return
	}

	respondWithJSON(w, http.StatusOK, match_players)
}
