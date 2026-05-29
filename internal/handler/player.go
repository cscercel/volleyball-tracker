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

type PlayerHandler struct {
	service	*service.PlayerService
}

func NewPlayerHandler(service *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{service: service}
}

func (h *PlayerHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/players", func(r chi.Router) {
		// Public routes
		r.Get("/", h.handleListPlayers)
		r.Get("/leaderboard", h.handleGetLeaderboard)
		r.Get("/{id}", h.handleGetPlayerByID)
		r.Get("/{id}/history", h.handleGetPlayerSeasonalMatches)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Post("/", h.handleCreatePlayer)
			r.Put("/{id}", h.handleUpdatePlayerName)
			r.Delete("/{id}", h.handleDeletePlayer)
		})
	})
}

// @Summary      Create Player
// @Tags         players
// @Produce      json
// @Security     BearerAuth
// @Param        body body      object{name=string} true "Player Body"
// @Success      201  {array}   db.Player
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/players [post]
func (h *PlayerHandler) handleCreatePlayer(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name	string 	`json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	player, err := h.service.CreatePlayer(r.Context(), body.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create player", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, player)
}

// @Summary      List Players
// @Tags         players
// @Produce      json
// @Success      200  {array}   db.Player
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/players [get]
func (h *PlayerHandler) handleListPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := h.service.ListPlayers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list players", err)
		return
	}

	respondWithJSON(w, http.StatusOK, players)
}

// @Summary      Update Player Name
// @Description  Updates the name of an existing player by their UUID
// @Tags         players
// @Accept       json
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string                true  "Player UUID" format(uuid)
// @Param        body body      object{name=string}   true  "New Player Name"
// @Success      200  {object}  db.Player
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/players/{id} [put]
func (h *PlayerHandler) handleUpdatePlayerName(w http.ResponseWriter, r *http.Request) {
	playerID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid player id", err)
		return
	}

	var body struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	player, err := h.service.UpdatePlayerName(r.Context(), playerID, body.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not change player name", err)
		return
	}

	respondWithJSON(w, http.StatusOK, player)
}

// @Summary      Delete Player
// @Tags         players
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string                true  "Player UUID" format(uuid)
// @Success      204  
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /api/v1/players/{id} [delete]
func (h *PlayerHandler) handleDeletePlayer(w http.ResponseWriter, r *http.Request) {
	playerID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid player id", err)
		return
	}

	if err := h.service.DeletePlayer(r.Context(), playerID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete character", err)
		return
	}
	
	respondWithJSON(w, http.StatusNoContent, "")
}

// @Summary      Get Player By ID
// @Tags         players
// @Produce      json
// @Param        id   		path  		string  true  "Player ID"
// @Param        match_type query  		string  true  "Match Type"
// @Param        season   	query      	integer true  "Season"
// @Success      200  		{array}   	db.GetPlayerStatsByIDRow
// @Failure      400  		{object}  	object{error=string}
// @Failure      404  		{object}  	object{error=string}
// @Router       /api/v1/players/{id} [get]
func (h *PlayerHandler) handleGetPlayerByID(w http.ResponseWriter, r *http.Request) {
	playerID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid player id", err)
		return
	}

	// Query params
	matchType := r.URL.Query().Get("match_type")
	if matchType != "indoor" && matchType != "beach" {
		respondWithError(w, http.StatusBadRequest, "invalid match type", 
			fmt.Errorf("expected: `indoor` or `beach`, got: %s", matchType),
		)
		return
	}

	seasonStr := r.URL.Query().Get("season")
	season, err := strconv.Atoi(seasonStr)
	if err != nil || season < 2023 {
		respondWithError(w, http.StatusBadRequest, "no season before 2023", err)
		return
	}

	player, err := h.service.GetPlayerByID(r.Context(), playerID, matchType, int32(season))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get player", err)
		return
	}

	respondWithJSON(w, http.StatusOK, player)
}

// @Summary      Get Leaderboard
// @Tags         players
// @Produce      json
// @Param        match_type query  		string  true  "Match Type"
// @Param        season   	query      	integer true  "Season"
// @Success      200  		{array}   	[]db.GetLeaderboardRow
// @Failure      400  		{object}  	object{error=string}
// @Failure      404  		{object}  	object{error=string}
// @Router       /api/v1/players/leaderboard [get]
func (h *PlayerHandler) handleGetLeaderboard(w http.ResponseWriter, r *http.Request) {
	// Query params
	matchType := r.URL.Query().Get("match_type")
	if matchType != "indoor" && matchType != "beach" {
		respondWithError(w, http.StatusBadRequest, "invalid match type", 
			fmt.Errorf("expected: `indoor` or `beach`, got: %s", matchType),
		)
		return
	}

	seasonStr := r.URL.Query().Get("season")
	season, err := strconv.Atoi(seasonStr)
	if err != nil || season < 2023 {
		respondWithError(w, http.StatusBadRequest, "no season before 2023", err)
		return
	}

	leaderboard, err := h.service.GetLeaderboard(r.Context(), matchType, int32(season))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to load leaderboard", err)
		return
	}

	respondWithJSON(w, http.StatusOK, leaderboard)
}

// @Summary      Get Player Match History
// @Tags         players
// @Produce      json
// @Param        id   		path  		string  true  "Player ID"
// @Param        match_type query  		string  true  "Match Type"
// @Param        season   	query      	integer true  "Season"
// @Success      200  		{array}   	[]db.GetPlayerSeasonalMatchesRow
// @Failure      400  		{object}  	object{error=string}
// @Failure      404  		{object}  	object{error=string}
// @Router       /api/v1/players/{id}/history [get]
func (h *PlayerHandler) handleGetPlayerSeasonalMatches(w http.ResponseWriter, r *http.Request) {
	playerID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid player id", err)
		return
	}

	// Query params
	matchType := r.URL.Query().Get("match_type")
	if matchType != "indoor" && matchType != "beach" {
		respondWithError(w, http.StatusBadRequest, "invalid match type", 
			fmt.Errorf("expected: `indoor` or `beach`, got: %s", matchType),
		)
		return
	}

	seasonStr := r.URL.Query().Get("season")
	season, err := strconv.Atoi(seasonStr)
	if err != nil || season < 2023 {
		respondWithError(w, http.StatusBadRequest, "no season before 2023", err)
		return
	}

	player_matches, err := h.service.GetPlayerSeasonalMatches(r.Context(), playerID, matchType, int32(season))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not load player history", err)
		return
	}

	respondWithJSON(w, http.StatusOK, player_matches)
}
