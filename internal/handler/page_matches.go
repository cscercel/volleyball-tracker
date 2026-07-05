package handler

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/cscercel/volleyball-tracker/web/templates/pages"
)

func (h *PageHandler) handleMatchesPage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	tab := q.Get("tab")
	if tab != "drafts" && tab != "completed" {
		tab = "create"
	}

	data := pages.MatchesPageData{
		ActiveTab: tab,
		LoggedIn:  isAuthenticated(r, h.jwtSecret),
	}

	switch tab {
	case "drafts":
		data.Drafts = h.buildDraftsTabData(r)
	case "completed":
		data.Completed = h.buildCompletedTabData(r, q)
	default:
		data.Create = h.buildCreateTabData(r, q)
	}

	if r.Header.Get("HX-Request") == "true" {
		pages.MatchesBody(data).Render(r.Context(), w)
		return
	}
	pages.Matches(data).Render(r.Context(), w)
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func joinCSV(names []string) string {
	return strings.Join(names, ",")
}

func contains(list []string, name string) bool {
	for _, n := range list {
		if n == name {
			return true
		}
	}
	return false
}

func removeName(list []string, name string) []string {
	out := make([]string, 0, len(list))
	for _, n := range list {
		if n != name {
			out = append(out, n)
		}
	}
	return out
}

// Court diagram
const (
	courtW = 450.0
	courtH = 270.0
	scaleX = courtW / 9
	scaleY = courtH / 18
)

var blueCoords = map[int][2]float64{
	1: {1, 2}, 2: {3.5, 2}, 3: {3.5, 9}, 4: {3.5, 16}, 5: {1, 16}, 6: {1, 9},
}
var redCoords = map[int][2]float64{
	7: {8, 16}, 8: {5.5, 16}, 9: {5.5, 9}, 10: {5.5, 2}, 11: {8, 2}, 12: {8, 9},
}
var positionSets = map[int][]int{
	1: {1},
	2: {1, 3},
	3: {1, 3, 5},
	4: {1, 2, 3, 5},
	5: {1, 2, 3, 4, 5},
	6: {1, 2, 3, 4, 5, 6},
}

func shuffleNames(names []string) []string {
	out := append([]string{}, names...)
	rand.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
	return out
}

func assignCourtSide(team []string, coords map[int][2]float64, offset int) ([]pages.CourtPlayer, []string) {
	shuffled := shuffleNames(team)
	onCourt := shuffled
	var bench []string
	if len(shuffled) > 6 {
		onCourt = shuffled[:6]
		bench = shuffled[6:]
	}

	positions := positionSets[len(onCourt)]
	players := make([]pages.CourtPlayer, len(onCourt))
	for i, name := range onCourt {
		pos := positions[i] + offset
		coord := coords[pos]
		players[i] = pages.CourtPlayer{Name: name, CX: coord[0] * scaleX, CY: coord[1] * scaleY}
	}
	return players, bench
}

func buildCourtData(blueTeam, redTeam []string) pages.CourtData {
	bluePlayers, blueBench := assignCourtSide(blueTeam, blueCoords, 0)
	redPlayers, redBench := assignCourtSide(redTeam, redCoords, 6)

	ballX, ballY := 0.2*scaleX, 2*scaleY
	if rand.Intn(2) == 1 {
		ballX, ballY = 8.8*scaleX, 16*scaleY
	}

	return pages.CourtData{
		BluePlayers: bluePlayers,
		RedPlayers:  redPlayers,
		BlueBench:   blueBench,
		RedBench:    redBench,
		BallX:       ballX,
		BallY:       ballY,
	}
}

func (h *PageHandler) buildCreateTabData(r *http.Request, q map[string][]string) pages.CreateTabData {
	get := func(key string) string {
		if v, ok := q[key]; ok && len(v) > 0 {
			return v[0]
		}
		return ""
	}

	matchType := get("match_type")
	if matchType != "indoor" && matchType != "beach" {
		matchType = "indoor"
	}
	draftType := get("draft_type")
	if draftType != "manual" {
		draftType = "random"
	}
	blueTeam := splitCSV(get("blue"))
	redTeam := splitCSV(get("red"))

	playerRows, _ := h.playerService.ListPlayers(r.Context())

	data := pages.CreateTabData{
		MatchType: matchType,
		DraftType: draftType,
		BlueTeam:  blueTeam,
		RedTeam:   redTeam,
		LoggedIn:  isAuthenticated(r, h.jwtSecret),
	}

	if len(playerRows) < 2 {
		data.NotEnoughPlayers = true
		return data
	}

	for _, p := range playerRows {
		if !contains(blueTeam, p.Name) && !contains(redTeam, p.Name) {
			data.Available = append(data.Available, pages.PlayerOption{ID: p.ID.String(), Name: p.Name})
		}
	}

	return data
}

func (h *PageHandler) handleCreateTeamUpdate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	blueTeam := splitCSV(r.FormValue("blue"))
	redTeam := splitCSV(r.FormValue("red"))

	if name := r.FormValue("add_blue"); name != "" && !contains(blueTeam, name) && !contains(redTeam, name) {
		blueTeam = append(blueTeam, name)
	}
	if name := r.FormValue("add_red"); name != "" && !contains(blueTeam, name) && !contains(redTeam, name) {
		redTeam = append(redTeam, name)
	}
	if name := r.FormValue("remove_blue"); name != "" {
		blueTeam = removeName(blueTeam, name)
	}
	if name := r.FormValue("remove_red"); name != "" {
		redTeam = removeName(redTeam, name)
	}

	q := map[string][]string{
		"match_type": {r.FormValue("match_type")},
		"draft_type": {r.FormValue("draft_type")},
		"blue":       {joinCSV(blueTeam)},
		"red":        {joinCSV(redTeam)},
	}
	data := h.buildCreateTabData(r, q)
	pages.CreateTab(data).Render(r.Context(), w)
}

func (h *PageHandler) handleCreateMatchSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	matchType := r.FormValue("match_type")
	draftType := r.FormValue("draft_type")
	blueTeam := splitCSV(r.FormValue("blue"))
	redTeam := splitCSV(r.FormValue("red"))

	if len(blueTeam) == 0 || len(redTeam) == 0 {
		data := h.buildCreateTabData(r, map[string][]string{
			"match_type": {matchType}, "draft_type": {draftType},
			"blue": {joinCSV(blueTeam)}, "red": {joinCSV(redTeam)},
		})
		data.Error = "Both teams need at least one player."
		pages.CreateTab(data).Render(r.Context(), w)
		return
	}

	finalBlue, finalRed := blueTeam, redTeam
	if draftType == "random" {
		all := append(append([]string{}, blueTeam...), redTeam...)
		rand.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })
		mid := len(all) / 2
		finalBlue, finalRed = all[:mid], all[mid:]
	}

	match, err := h.matchService.CreateMatch(r.Context(), matchType, finalBlue, finalRed)
	if err != nil {
		data := h.buildCreateTabData(r, map[string][]string{"match_type": {matchType}, "draft_type": {draftType}})
		data.Error = "Failed to create match."
		pages.CreateTab(data).Render(r.Context(), w)
		return
	}

	roster, _ := h.matchService.GetMatchPlayers(r.Context(), match.ID)
	data := h.buildCreateTabData(r, map[string][]string{"match_type": {matchType}, "draft_type": {draftType}})
	data.Success = "✅ Match created!"
	for _, p := range roster {
		if p.Color == "blue" {
			data.ResultBlue = append(data.ResultBlue, p.PlayerName)
		} else {
			data.ResultRed = append(data.ResultRed, p.PlayerName)
		}
	}
	if len(data.ResultBlue) > 0 || len(data.ResultRed) > 0 {
		data.Court = buildCourtData(data.ResultBlue, data.ResultRed)
		data.HasCourt = true
	}
	pages.CreateTab(data).Render(r.Context(), w)
}

func (h *PageHandler) buildDraftsTabData(r *http.Request) pages.DraftsTabData {
	matches, err := h.matchService.ListUncompletedMatches(r.Context())
	if err != nil {
		return pages.DraftsTabData{
			Error: "Failed to load draft matches.", 
			LoggedIn: isAuthenticated(r, h.jwtSecret),
		}
	}

	data := pages.DraftsTabData{LoggedIn: isAuthenticated(r, h.jwtSecret)}
	for _, m := range matches {
		roster, _ := h.matchService.GetMatchPlayers(r.Context(), m.ID)
		draft := pages.DraftMatch{
			ID:        m.ID.String(),
			MatchType: m.MatchType,
		}
		for _, p := range roster {
			if p.Color == "blue" {
				draft.BlueRoster = append(draft.BlueRoster, p.PlayerName)
			} else {
				draft.RedRoster = append(draft.RedRoster, p.PlayerName)
			}
		}
		data.Matches = append(data.Matches, draft)
	}
	return data
}

func (h *PageHandler) handleDraftSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		pages.DraftsTab(pages.DraftsTabData{Error: "Bad form data."}).Render(r.Context(), w)
		return
	}

	matchID, err := uuid.Parse(r.FormValue("match_id"))
	if err != nil {
		pages.DraftsTab(pages.DraftsTabData{Error: "Invalid match."}).Render(r.Context(), w)
		return
	}
	blueScore, errB := strconv.Atoi(r.FormValue("blue_score"))
	redScore, errR := strconv.Atoi(r.FormValue("red_score"))
	if errB != nil || errR != nil {
		data := h.buildDraftsTabData(r)
		data.Error = "Enter scores for both teams."
		pages.DraftsTab(data).Render(r.Context(), w)
		return
	}
	if blueScore == redScore {
		data := h.buildDraftsTabData(r)
		data.Error = "Scores cannot be equal."
		pages.DraftsTab(data).Render(r.Context(), w)
		return
	}

	if _, err := h.matchService.RegisterMatch(r.Context(), matchID, int32(blueScore), int32(redScore)); err != nil {
		data := h.buildDraftsTabData(r)
		data.Error = "Failed to submit results."
		pages.DraftsTab(data).Render(r.Context(), w)
		return
	}

	pages.DraftsTab(h.buildDraftsTabData(r)).Render(r.Context(), w)
}

func (h *PageHandler) handleDraftDelete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		pages.DraftsTab(pages.DraftsTabData{Error: "Bad form data."}).Render(r.Context(), w)
		return
	}

	matchID, err := uuid.Parse(r.FormValue("match_id"))
	if err != nil {
		pages.DraftsTab(pages.DraftsTabData{Error: "Invalid match."}).Render(r.Context(), w)
		return
	}

	if err := h.matchService.DeleteUncompletedMatch(r.Context(), matchID); err != nil {
		data := h.buildDraftsTabData(r)
		data.Error = "Failed to delete draft."
		pages.DraftsTab(data).Render(r.Context(), w)
		return
	}

	pages.DraftsTab(h.buildDraftsTabData(r)).Render(r.Context(), w)
}

func (h *PageHandler) buildCompletedTabData(r *http.Request, q map[string][]string) pages.CompletedTabData {
	get := func(key string) string {
		if v, ok := q[key]; ok && len(v) > 0 {
			return v[0]
		}
		return ""
	}

	matchType := get("match_type")
	if matchType != "indoor" && matchType != "beach" {
		matchType = "indoor"
	}
	season, err := strconv.Atoi(get("season"))
	if err != nil || season < 2023 {
		season = time.Now().Year()
	}

	data := pages.CompletedTabData{MatchType: matchType, Season: season}

	matches, err := h.matchService.ListMatchesBySeason(r.Context(), matchType, int32(season))
	if err != nil {
		data.Error = "Failed to load completed matches."
		return data
	}

	for _, m := range matches {
		roster, _ := h.matchService.GetMatchPlayers(r.Context(), m.ID)
		diff := m.BlueScore - m.RedScore
		if diff < 0 {
			diff = -diff
		}
		winner := "red"
		if m.BlueScore > m.RedScore {
			winner = "blue"
		}
		date := ""
		if m.CreatedAt.Valid {
			date = m.CreatedAt.Time.Format("2006-01-02")
		}

		cm := pages.CompletedMatch{
			ID:        m.ID.String(),
			MatchType: m.MatchType,
			BlueScore: m.BlueScore,
			RedScore:  m.RedScore,
			IsOtl:     diff == 2,
			Winner:    winner,
			Date:      date,
		}
		for _, p := range roster {
			if p.Color == "blue" {
				cm.BlueRoster = append(cm.BlueRoster, p.PlayerName)
			} else {
				cm.RedRoster = append(cm.RedRoster, p.PlayerName)
			}
		}
		data.Matches = append(data.Matches, cm)
	}

	return data
}
