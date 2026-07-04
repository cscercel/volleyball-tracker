package handler

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/cscercel/volleyball-tracker/web/templates/pages"
)

type leaderboardRow struct {
	Name           string
	Played         int32
	Wins           int32
	Losses         int32
	Otl            int32
	Points         int32
	WinRate        float64
	EfficiencyRate float64
}

func toFloat64(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

func (h *PageHandler) handleHomePage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	matchType := q.Get("match_type")
	if matchType != "indoor" && matchType != "beach" {
		matchType = "indoor"
	}

	currentYear := time.Now().Year()
	season, err := strconv.Atoi(q.Get("season"))
	if err != nil || season < 2023 {
		season = currentYear
	}

	sortColumn := q.Get("sort")
	if sortColumn == "" {
		sortColumn = "points"
	}
	sortDir := q.Get("dir")
	if sortDir != "asc" {
		sortDir = "desc"
	}

	rows, err := h.playerService.GetLeaderboard(r.Context(), matchType, int32(season))

	var leaderboard []leaderboardRow
	loadErr := ""
	if err != nil {
		loadErr = "Failed to load leaderboard."
	} else {
		leaderboard = make([]leaderboardRow, 0, len(rows))
		for _, row := range rows {
			winRate := toFloat64(row.WinRate)
			efficiencyRate := toFloat64(row.EfficiencyRate)
			leaderboard = append(leaderboard, leaderboardRow{
				Name:           row.Name,
				Played:         row.Played,
				Wins:           row.Wins,
				Losses:         row.Losses,
				Otl:            row.Otl,
				Points:         row.Points,
				WinRate:        winRate,
				EfficiencyRate: efficiencyRate,
			})
		}
		sortLeaderboard(leaderboard, sortColumn, sortDir)
	}

	data := pages.LeaderboardData{
		MatchType:  matchType,
		Season:     season,
		SortColumn: sortColumn,
		SortDir:    sortDir,
		Rows:       toTemplRows(leaderboard),
		Error:      loadErr,
		LoggedIn:   isAuthenticated(r, h.jwtSecret),
	}

	if r.Header.Get("HX-Request") == "true" {
		pages.LeaderboardTable(data).Render(r.Context(), w)
		return
	}

	pages.Home(data).Render(r.Context(), w)
}

func toTemplRows(rows []leaderboardRow) []pages.LeaderboardRow {
	out := make([]pages.LeaderboardRow, len(rows))
	for i, row := range rows {
		out[i] = pages.LeaderboardRow{
			Rank:    i + 1,
			Name:    row.Name,
			Played:  row.Played,
			Wins:    row.Wins,
			Losses:  row.Losses,
			Otl:     row.Otl,
			Points:  row.Points,
			WinRate: row.WinRate,
		}
	}
	return out
}

func sortLeaderboard(rows []leaderboardRow, column, dir string) {
	less := func(i, j int) bool {
		var result bool
		switch column {
		case "name":
			result = rows[i].Name < rows[j].Name
		case "played":
			result = rows[i].Played < rows[j].Played
		case "wins":
			result = rows[i].Wins < rows[j].Wins
		case "losses":
			result = rows[i].Losses < rows[j].Losses
		case "otl":
			result = rows[i].Otl < rows[j].Otl
		case "win_rate":
			result = rows[i].WinRate < rows[j].WinRate
		default:
			result = rows[i].Points < rows[j].Points
		}
		if dir == "desc" {
			return !result
		}
		return result
	}
	sort.SliceStable(rows, less)
}
