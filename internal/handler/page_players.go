package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/cscercel/volleyball-tracker/internal/db"
	"github.com/cscercel/volleyball-tracker/web/templates/pages"
)

func (h *PageHandler) handlePlayersPage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	tab := q.Get("tab")
	if tab != "add" && tab != "manage" {
		tab = "profile"
	}

	playerRows, err := h.playerService.ListPlayers(r.Context())
	if err != nil {
		log.Printf("handlePlayersPage: ListPlayers failed: %v", err)
	}
	players := make([]pages.PlayerOption, 0, len(playerRows))
	for _, p := range playerRows {
		players = append(players, pages.PlayerOption{ID: p.ID.String(), Name: p.Name})
	}

	data := pages.PlayersPageData{
		ActiveTab: tab,
		Players:   players,
	}

	switch tab {
	case "profile":
		data.Profile = h.buildProfileTabData(r, q, players)
	case "manage":
		data.Manage = pages.ManageTabData{
			Players:          players,
			SelectedPlayerID: q.Get("player_id"),
		}
	case "add":
	}

	if r.Header.Get("HX-Request") == "true" {
		pages.PlayersTabContent(data).Render(r.Context(), w)
		return
	}
	pages.Players(data).Render(r.Context(), w)
}

func (h *PageHandler) buildProfileTabData(
	r *http.Request, q map[string][]string, players []pages.PlayerOption,
) pages.ProfileData {
	get := func(key string) string {
		if v, ok := q[key]; ok && len(v) > 0 {
			return v[0]
		}
		return ""
	}

	selectedID := get("player_id")
	matchType := get("match_type")
	if matchType != "indoor" && matchType != "beach" {
		matchType = "indoor"
	}
	season, err := strconv.Atoi(get("season"))
	if err != nil || season < 2023 {
		season = time.Now().Year()
	}

	profile := pages.ProfileData{
		Players:          players,
		SelectedPlayerID: selectedID,
		MatchType:        matchType,
		Season:           season,
	}

	if selectedID == "" {
		return profile
	}

	playerID, err := uuid.Parse(selectedID)
	if err != nil {
		profile.Error = "Invalid player."
		return profile
	}

	stats, err := h.playerService.GetPlayerByID(r.Context(), playerID, matchType, int32(season))
	if err != nil {
		profile.Error = "Failed to load player stats."
		return profile
	}

	history, _ := h.playerService.GetPlayerSeasonalMatches(r.Context(), playerID, matchType, int32(season))

	prevStats, prevErr := h.playerService.GetPlayerByID(r.Context(), playerID, matchType, int32(season-1))
	hasPrev := prevErr == nil

	winRate := toFloat64(stats.WinRate)

	prevWinRatePct := 0
	if hasPrev {
		prevWinRatePct = int(toFloat64(prevStats.WinRate) * 100)
	}

	profile.HasStats = true
	profile.Stats = []pages.StatCard{
		statCard("Matches Played", int(stats.Played), int(prevStats.Played), hasPrev, ""),
		statCard("Wins", int(stats.Wins), int(prevStats.Wins), hasPrev, ""),
		statCard("Losses", int(stats.Losses), int(prevStats.Losses), hasPrev, ""),
		statCard("OTL", int(stats.Otl), int(prevStats.Otl), hasPrev, ""),
		statCard("Points", int(stats.Points), int(prevStats.Points), hasPrev, ""),
		statCard("Win Rate", int(winRate*100), prevWinRatePct, hasPrev, "%"),
		statCard("Win Streak", int(stats.Streak), 0, false, ""), // old UI never showed a delta for streak
		statCard("Longest Streak", int(stats.LongestStreak), int(prevStats.LongestStreak), hasPrev, ""),
	}

	profile.History = make([]pages.HistoryRow, 0, len(history))
	for _, m := range history {
		profile.History = append(profile.History, buildHistoryRow(m))
	}

	return profile
}

func statCard(label string, value, prev int, hasPrev bool, suffix string) pages.StatCard {
	card := pages.StatCard{
		Label: label,
		Value: strconv.Itoa(value) + suffix,
	}
	if !hasPrev {
		return card
	}
	diff := value - prev
	switch {
	case diff == 0:
		card.Delta = "±0"
		card.DeltaClass = "neutral"
	case diff > 0:
		card.Delta = "+" + strconv.Itoa(diff)
		card.DeltaClass = "positive"
	default:
		card.Delta = strconv.Itoa(diff)
		card.DeltaClass = "negative"
	}
	return card
}

func buildHistoryRow(m db.GetPlayerSeasonalMatchesRow) pages.HistoryRow {
	myScore, theirScore := m.RedScore, m.BlueScore
	if m.Color == "blue" {
		myScore, theirScore = m.BlueScore, m.RedScore
	}
	won := myScore > theirScore
	diff := m.BlueScore - m.RedScore
	if diff < 0 {
		diff = -diff
	}
	isOtl := diff == 2 && !won

	result, class := "Loss", "loss"
	switch {
	case won:
		result, class = "Win", "win"
	case isOtl:
		result, class = "OTL", "otl"
	}

	date := ""
	if m.CreatedAt.Valid {
		date = m.CreatedAt.Time.Format("2006-01-02")
	}

	return pages.HistoryRow{
		Result:      result,
		ResultClass: class,
		Score:       strconv.Itoa(int(myScore)) + " : " + strconv.Itoa(int(theirScore)),
		Team:        m.Color + " team",
		Date:        date,
	}
}

func (h *PageHandler) handleAddPlayerSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		pages.AddTab(pages.AddTabData{Error: "Could not read form data."}).Render(r.Context(), w)
		return
	}

	name := r.FormValue("name")
	if _, err := h.playerService.CreatePlayer(r.Context(), name); err != nil {
		pages.AddTab(pages.AddTabData{Error: "Failed to add player."}).Render(r.Context(), w)
		return
	}

	pages.AddTab(pages.AddTabData{Success: "Added " + name}).Render(r.Context(), w)
}

func (h *PageHandler) handleRenamePlayerSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		pages.ManageTab(pages.ManageTabData{Error: "Could not read form data."}).Render(r.Context(), w)
		return
	}
	h.submitManageAction(w, r, func(playerID uuid.UUID) error {
		_, err := h.playerService.UpdatePlayerName(r.Context(), playerID, r.FormValue("new_name"))
		return err
	}, "Renamed successfully", "Failed to rename player.", true)
}

func (h *PageHandler) handleDeletePlayerSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		pages.ManageTab(pages.ManageTabData{Error: "Could not read form data."}).Render(r.Context(), w)
		return
	}
	h.submitManageAction(w, r, func(playerID uuid.UUID) error {
		return h.playerService.DeletePlayer(r.Context(), playerID)
	}, "Player deleted", "Failed to delete player.", false)
}

func (h *PageHandler) submitManageAction(
	w http.ResponseWriter, r *http.Request,
	action func(uuid.UUID) error,
	successMsg, errorMsg string,
	keepSelected bool,
) {
	playerIDStr := r.FormValue("player_id")
	playerID, err := uuid.Parse(playerIDStr)
	if err != nil {
		pages.ManageTab(pages.ManageTabData{Error: "Invalid player."}).Render(r.Context(), w)
		return
	}

	data := pages.ManageTabData{}
	if err := action(playerID); err != nil {
		data.Error = errorMsg
		data.SelectedPlayerID = playerIDStr
	} else {
		data.Success = successMsg
		if keepSelected {
			data.SelectedPlayerID = playerIDStr
		}
	}

	rows, _ := h.playerService.ListPlayers(r.Context())
	data.Players = make([]pages.PlayerOption, 0, len(rows))
	for _, p := range rows {
		data.Players = append(data.Players, pages.PlayerOption{ID: p.ID.String(), Name: p.Name})
	}

	pages.ManageTab(data).Render(r.Context(), w)
}
