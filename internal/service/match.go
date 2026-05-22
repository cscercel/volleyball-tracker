package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/cscercel/volleyball-tracker/internal/db"
	"github.com/google/uuid"
)

type MatchService struct {
	queries	*db.Queries
}

type TeamPlayer struct {
	PlayerID	uuid.UUID	`json:"player_id"`
	Name		string		`json:"name"`
}

type MatchWithPlayers struct {
	MatchID		uuid.UUID			`json:"match_id"`
	MatchType	string				`json:"match_type"`
	Season		int					`json:"season"`
	BlueScore	int					`json:"blue_score"`
	RedScore	int					`json:"red_score"`
	IsCompleted	bool				`json:"is_completed"`
	BlueTeam	[]TeamPlayer		`json:"blue_team"`
	RedTeam		[]TeamPlayer		`json:"red_team"`
}

type TeamPerformance struct {
	MatchType	string				`json:"match_type"`
	Season		int					`json:"season"`
	Players		[]TeamPlayer		`json:"players"`
	Scored		int					`json:"scored"`
	Conceded	int					`json:"conceded"`
	IsWinner	bool				`json:"is_winner"`
	IsOtl		bool				`json:"is_otl"`
}

func NewMatchService(queries *db.Queries) *MatchService {
	return &MatchService{queries: queries}
}

// Helper method to split match players into respective teams
func (s *MatchService) groupIntoTeam(
	ctx context.Context, match_players []db.MatchPlayer, color string,
) ([]TeamPlayer, error) {
	if color != "blue" && color != "red" {
		return []TeamPlayer{}, fmt.Errorf("color must be either blue or red, got: %s", color)
	}

	team := []TeamPlayer{}

	for _, match_player := range match_players {
		player, err := s.queries.GetPlayer(ctx, match_player.PlayerID)
		if err != nil {
			return []TeamPlayer{}, fmt.Errorf("failed to load player: %w", err)
		}

		if match_player.Color == color {
			team = append(team, TeamPlayer{
				PlayerID: player.ID, 
				Name: player.Name, 
			})
		}
	}

	return team, nil
}

func (s *MatchService) GetMatch(ctx context.Context, match_id uuid.UUID) (MatchWithPlayers, error) {
	match, err := s.queries.GetMatch(ctx, match_id)
	if err != nil {
		return MatchWithPlayers{}, fmt.Errorf("unable to get match: %w", err)
	}

	players, err := s.queries.GetPlayersFromMatch(ctx, match.ID)
	if err != nil {
		return MatchWithPlayers{}, fmt.Errorf("failed to load match players: %w", err)
	}

	blue_team, err := s.groupIntoTeam(ctx, players, "blue")
	if err != nil {
		return MatchWithPlayers{}, fmt.Errorf("failed to load blue team: %w", err)
	}

	red_team, err := s.groupIntoTeam(ctx, players, "red")
	if err != nil {
		return MatchWithPlayers{}, fmt.Errorf("failed to load red team: %w", err)
	}

	return MatchWithPlayers{
		MatchID: match.ID,
		MatchType: match.MatchType,
		Season: int(match.Season),
		BlueScore: int(match.BlueScore),
		RedScore: int(match.RedScore),
		IsCompleted: match.IsCompleted,
		BlueTeam: blue_team,
		RedTeam: red_team,
	}, nil
}

func (s *MatchService) GetRegisteredMatches(ctx context.Context) ([]MatchWithPlayers, error) {
	matches, err := s.queries.GetRegisteredMatches(ctx)
	if err != nil {
		return []MatchWithPlayers{}, fmt.Errorf("unable to get registered matches: %w", err)
	}

	matches_with_players := []MatchWithPlayers{}
	for _, match := range matches {
		blue_team, err := s.queries.GetBlueTeamFromMatch(ctx, match.ID)
		if err != nil {
			return []MatchWithPlayers{}, fmt.Errorf("failed to load blue team: %w", err)
		}

		red_team, err := s.queries.GetRedTeamFromMatch(ctx, match.ID)
		if err != nil {
			return []MatchWithPlayers{}, fmt.Errorf("failed to load red team: %w", err)
		}

		matches_with_players = append(matches_with_players, MatchWithPlayers{
			Match: match,
			BlueTeam: blue_team,
			RedTeam: red_team,
		})
	}

	return matches_with_players, nil
}

func (s *MatchService) GetDrafts(ctx context.Context) ([]MatchWithPlayers, error) {
	drafts, err := s.queries.GetDrafts(ctx)
	if err != nil {
		return []MatchWithPlayers{}, fmt.Errorf("unable to get drafts: %w", err)
	}

	matches_with_players := []MatchWithPlayers{}
	for _, match := range drafts {
		blue_team, err := s.queries.GetBlueTeamFromMatch(ctx, match.ID)
		if err != nil {
			return []MatchWithPlayers{}, fmt.Errorf("failed to load blue team: %w", err)
		}

		red_team, err := s.queries.GetRedTeamFromMatch(ctx, match.ID)
		if err != nil {
			return []MatchWithPlayers{}, fmt.Errorf("failed to load red team: %w", err)
		}

		matches_with_players = append(matches_with_players, MatchWithPlayers{
			Match: match,
			BlueTeam: blue_team,
			RedTeam: red_team,
		})
	}

	return matches_with_players, nil
}

func (s *MatchService) GetSeasonMatches(
	ctx context.Context, match_type string, season int,
) ([]MatchWithPlayers, error) {
	params := db.GetSeasonMatchesParams{
		MatchType: match_type,
		Season: int32(season),
	}

	matches, err := s.queries.GetSeasonMatches(ctx, params)
	if err != nil {
		return []MatchWithPlayers{}, fmt.Errorf("unable to get seasonal matches: %w", err)
	}

	matches_with_players := []MatchWithPlayers{}
	for _, match := range matches {
		blue_team, err := s.queries.GetBlueTeamFromMatch(ctx, match.ID)
		if err != nil {
			return []MatchWithPlayers{}, fmt.Errorf("failed to load blue team: %w", err)
		}

		red_team, err := s.queries.GetRedTeamFromMatch(ctx, match.ID)
		if err != nil {
			return []MatchWithPlayers{}, fmt.Errorf("failed to load red team: %w", err)
		}

		matches_with_players = append(matches_with_players, MatchWithPlayers{
			Match: match,
			BlueTeam: blue_team,
			RedTeam: red_team,
		})
	}

	return matches_with_players, nil
}

func (s *MatchService) CreateMatch(
	ctx context.Context, match_type string, blue_players, red_players []uuid.UUID,
) (MatchWithPlayers, error) {
	if match_type != "indoor" && match_type != "beach" {
		return MatchWithPlayers{}, fmt.Errorf("invalid match_type: expected `indoor` or `beach`, got %s", match_type)
	}

	if len(blue_players) == 0 || len(red_players) == 0 {
		return MatchWithPlayers{}, fmt.Errorf("both teams must have at least 1 player")
	}

	season := time.Now().UTC().Year()

	params := db.CreateMatchParams{
		MatchType: match_type,
		Season: int32(season),
	}

	match, err := s.queries.CreateMatch(ctx, params)
	if err != nil {
		return MatchWithPlayers{}, fmt.Errorf("could not create match: %w", err)
	}

	// Add players from blue_team
	blue_team := []db.MatchPlayer{}
	for _, player := range blue_players {
		player, err := s.queries.AddPlayerToMatch(ctx, db.AddPlayerToMatchParams{
			MatchID: match.ID,
			PlayerID: player,
			Color: "blue",
		})
		if err != nil {
			return MatchWithPlayers{}, fmt.Errorf("failed to add players from blue_team: %w", err)
		}

		blue_team = append(blue_team, player)
	}

	// idem. for red_team
	red_team := []db.MatchPlayer{}
	for _, player := range red_players {
		player, err := s.queries.AddPlayerToMatch(ctx, db.AddPlayerToMatchParams{
			MatchID: match.ID,
			PlayerID: player,
			Color: "red",
		})
		if err != nil {
			return MatchWithPlayers{}, fmt.Errorf("failed to add players from red_team: %w", err)
		}

		red_team = append(red_team, player)
	}

	return MatchWithPlayers{
		Match: match,
		BlueTeam: blue_team,
		RedTeam: red_team,
	}, nil
}

// Helper function for updating player stats after match registration
func (s *MatchService) updateTeamStats(ctx context.Context, team TeamPerformance) error {
	switch {
	case team.IsWinner:
		for _, player := range team.Players {
			_, err := s.queries.UpdatePlayerStatsWin(ctx, db.UpdatePlayerStatsWinParams{
				PlayerID: player.PlayerID,
				MatchType: team.MatchType,
				Season: int32(team.Season),
				Scored: int32(team.Scored),
				Conceded: int32(team.Conceded),	
			})
			if err != nil {
				return fmt.Errorf("failed to declare player as winner: %w", err)
			}
		}
	case team.IsOtl:
		for _, player := range team.Players {
			_, err := s.queries.UpdatePlayerStatsOtl(ctx, db.UpdatePlayerStatsOtlParams{
				PlayerID: player.PlayerID,
				MatchType: team.MatchType,
				Season: int32(team.Season),
				Scored: int32(team.Scored),
				Conceded: int32(team.Conceded),	
			})
			if err != nil {
				return fmt.Errorf("failed to declare player as winner: %w", err)
			}
		}
	default:
		for _, player := range team.Players {
			_, err := s.queries.UpdatePlayerStatsLoss(ctx, db.UpdatePlayerStatsLossParams{
				PlayerID: player.PlayerID,
				MatchType: team.MatchType,
				Season: int32(team.Season),
				Scored: int32(team.Scored),
				Conceded: int32(team.Conceded),	
			})
			if err != nil {
				return fmt.Errorf("failed to declare player as winner: %w", err)
			}
		}
	}

	return nil
}

func (s *MatchService) RegisterMatch(
	ctx context.Context, match_id uuid.UUID, blue_score, red_score int,
) (MatchWithPlayers, error) {
	// Check scores
	if blue_score == red_score {
		return MatchWithPlayers{}, fmt.Errorf("cannot determine winner")
	}

	is_overtime := false
	if math.Abs(float64(blue_score) - float64(red_score)) == 2 {
		is_overtime = true
	}

	// Get Match
	params := db.RegisterMatchParams{
		ID: match_id,
		BlueScore: int32(blue_score),
		RedScore: int32(red_score),
	}
	match, err := s.queries.RegisterMatch(ctx, params)	
	if err != nil {
		return MatchWithPlayers{}, fmt.Errorf("failed to get match for registration: %w", err)
	}

	// Get Blue Team & Red Team
	blue_players, err := s.queries.GetBlueTeamFromMatch(ctx, match.ID)
	if err != nil {
		return MatchWithPlayers{}, fmt.Errorf("failed to load blue team from match: %w", err)
	}

	blue_team := TeamPerformance{
		MatchType: match.MatchType,
		Season: int(match.Season),
		Players: blue_players,
		Scored: blue_score,
		Conceded: red_score,
		IsWinner: (blue_score > red_score),
		IsOtl: (blue_score < red_score) && is_overtime,
	}

	red_players, err := s.queries.GetRedTeamFromMatch(ctx, match.ID)
	if err != nil {
		return MatchWithPlayers{}, fmt.Errorf("failed to load red team from match: %w", err)
	}

	red_team := TeamPerformance{
		MatchType: match.MatchType,
		Season: int(match.Season),
		Players: red_players,
		Scored: red_score,
		Conceded: blue_score,
		IsWinner: (red_score > blue_score),
		IsOtl: (red_score < blue_score) && is_overtime,
	}

	// Match results
	if err := s.updateTeamStats(ctx, blue_team); err != nil {
		return MatchWithPlayers{}, fmt.Errorf("failed to register match for blue team: %w", err)
	}

	if err := s.updateTeamStats(ctx, red_team); err != nil {
		return MatchWithPlayers{}, fmt.Errorf("failed to register match for red team: %w", err)
	}

	return MatchWithPlayers{
		Match: match,
		BlueTeam: blue_players,
		RedTeam: red_players,
	}, nil
}

func (s *MatchService) DeleteDraft(ctx context.Context, match_id uuid.UUID) error {
	// Check if match was completed
	match, err := s.queries.GetMatch(ctx, match_id)
	if err != nil {
		return fmt.Errorf("failed to get match: %w", err)
	}

	if !match.IsCompleted {
		return fmt.Errorf("cannot delete a registered match")
	}

	if err := s.queries.DeleteDraft(ctx, match_id); err != nil {
		return fmt.Errorf("could not delete match: %w", err)
	}
	
	return nil
}
