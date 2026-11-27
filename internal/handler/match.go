package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cscercel/volleyball-tracker/internal/db"
	"github.com/cscercel/volleyball-tracker/internal/service"
	"github.com/go-chi/chi/v5"
)

type MatchHandler struct {
	service *service.MatchService
}

func NewMatchHandler(service *service.MatchService) *MatchHandler {
	return &MatchHandler{service: service}
}

func (h *MatchHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/matches", func(r chi.Router) {
		// Public routes
		r.Get("/{id}", h.handleGetMatch)
		r.Get("/", h.handleListMatchesBySeason)
		r.Get("/uncompleted", h.handleListUncompletedMatches)
		r.Get("/{id}/roster", h.handleGetMatchPlayers)
		r.Post("/demo", h.handleCreateDemoMatch)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Post("/", h.handleCreateMatch)
			r.Delete("/{id}", h.handleDeleteUncompletedMatch)
			r.Put("/{id}", h.handleRegisterMatch)
		})
	})
}

// @Summary      Create Match
// @Tags         matches
// @Produce      json
// @Security     BearerAuth
// @Param        body body      object{match_type=string, blue_team=[]string, red_team=[]string} true "Match Body"
// @Success      201  {array}   db.Match
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/matches [post]
func (h *MatchHandler) handleCreateMatch(w http.ResponseWriter, r *http.Request) {
	var body struct {
		MatchType string   `json:"match_type"`
		BlueTeam  []string `json:"blue_team"`
		RedTeam   []string `json:"red_team"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	match, err := h.service.CreateMatch(r.Context(), body.MatchType, body.BlueTeam, body.RedTeam, false)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create match", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, match)
}

// @Summary      Create Demo Match
// @Tags         matches
// @Produce      json
// @Param        body body      object{match_type=string, blue_team=[]string, red_team=[]string} true "Match Body"
// @Success      201  {array}   db.Match
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/demo [post]
func (h *MatchHandler) handleCreateDemoMatch(w http.ResponseWriter, r *http.Request) {
	var body struct {
		MatchType string   `json:"match_type"`
		BlueTeam  []string `json:"blue_team"`
		RedTeam   []string `json:"red_team"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	match, err := h.service.CreateMatch(r.Context(), body.MatchType, body.BlueTeam, body.RedTeam, true)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create match", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, match)
}

// @Summary      Get Match
// @Tags         matches
// @Produce      json
// @Param        id   path      int                true  "Match ID"
// @Success      200  {array}   db.Match
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/matches/{id} [get]
func (h *MatchHandler) handleGetMatch(w http.ResponseWriter, r *http.Request) {
	matchIDstr := chi.URLParam(r, "id")

	matchID, err := strconv.ParseInt(matchIDstr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "id must be an integer", err)
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
// @Param        page   	query      	integer false  "Page number"
// @Param        page_limit query      	integer false  "Page limit"
// @Success      200  		{object}  	object{matches=[]db.Match,has_more=bool}
// @Failure      400  		{object}  	object{error=string}
// @Failure      500  		{object}  	object{error=string}
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
	season, err := strconv.ParseInt(seasonStr, 10, 64)
	if err != nil || season < 2023 {
		respondWithError(w, http.StatusBadRequest, "no season before 2023", err)
		return
	}

	page := int64(1)
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil || page < 1 {
			respondWithError(w, http.StatusBadRequest, "invalid page", err)
			return
		}
	}

	pageLimit := int64(5)
	if pageLimitStr := r.URL.Query().Get("page_limit"); pageLimitStr != "" {
		pageLimit, err = strconv.ParseInt(pageLimitStr, 10, 64)
		if err != nil || pageLimit < 1 {
			respondWithError(w, http.StatusBadRequest, "invalid page limit", err)
			return
		}
	}

	matches, hasMore, err := h.service.ListMatchesBySeason(
		r.Context(), matchType, season, page, pageLimit,
	)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get matches", err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Matches []db.Match `json:"matches"`
		HasMore bool       `json:"has_more"`
	}{
		Matches: matches,
		HasMore: hasMore,
	})
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
// @Security     BearerAuth
// @Param        id   path      int                true  "Match ID"
// @Success      204
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/matches/{id} [delete]
func (h *MatchHandler) handleDeleteUncompletedMatch(w http.ResponseWriter, r *http.Request) {
	matchIDstr := chi.URLParam(r, "id")

	matchID, err := strconv.ParseInt(matchIDstr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "id must be an integer", err)
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
// @Param        id   path      int                true  "Match ID"
// @Success      200  {array}   []db.GetMatchPlayersRow
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/matches/{id}/roster [get]
func (h *MatchHandler) handleGetMatchPlayers(w http.ResponseWriter, r *http.Request) {
	matchIDstr := chi.URLParam(r, "id")

	matchID, err := strconv.ParseInt(matchIDstr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "id must be an integer", err)
		return
	}

	match_players, err := h.service.GetMatchPlayers(r.Context(), matchID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get match players", err)
		return
	}

	respondWithJSON(w, http.StatusOK, match_players)
}

// @Summary      Update Match Scores
// @Description  Updates the scores of a match
// @Tags         matches
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int               true  "Match ID"
// @Param        body body      object{blue_score=int, red_score=int}   true  "Match Scores"
// @Success      200  {object}  db.Match
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/matches/{id} [put]
func (h *MatchHandler) handleRegisterMatch(w http.ResponseWriter, r *http.Request) {
	matchIDstr := chi.URLParam(r, "id")

	matchID, err := strconv.ParseInt(matchIDstr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "id must be an integer", err)
		return
	}

	var body struct {
		BlueScore int64 `json:"blue_score"`
		RedScore  int64 `json:"red_score"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid body type", err)
		return
	}

	match, err := h.service.RegisterMatch(r.Context(), matchID, body.BlueScore, body.RedScore)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not register match", err)
		return
	}

	respondWithJSON(w, http.StatusOK, match)
}
