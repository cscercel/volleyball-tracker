package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cscercel/volleyball-tracker/internal/db"
	"github.com/google/uuid"
)

type MatchService struct {
	queries	*db.Queries
}

func NewMatchService(queries *db.Queries) *MatchService {
	return &MatchService{queries: queries}
}

// Helper function to add players to match
func (s *MatchService) addPlayerToMatch(
	ctx context.Context, match db.Match, player db.Player, color string,
) error {
	if color != "blue" && color != "red" {
		return fmt.Errorf("incorrect color, expected `blue` or `red`, got `%s`", color)
	}

	// Check if player has stats for match type & season
	_, err := s.queries.GetPlayerStatsByID(ctx, db.GetPlayerStatsByIDParams{
		PlayerID: player.ID,
		MatchType: match.MatchType,
		Season: match.Season,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = s.queries.UpsertPlayerStats(ctx, db.UpsertPlayerStatsParams{
				PlayerID: player.ID,
				MatchType: match.MatchType,
				Season: match.Season,
			})
			if err != nil {
				return fmt.Errorf("failed to initliaze stats for player: %w", err)
			}
		}
		return fmt.Errorf("failed to load stats for player: %w", err)
	}

	_, err = s.queries.AddPlayerToMatch(ctx, db.AddPlayerToMatchParams{
		MatchID: match.ID,
		PlayerID: player.ID,
		Color: color,
	})
	if err != nil {
		return fmt.Errorf("failed to add player to match: %w", err)
	}

	return nil
}

func (s *MatchService) CreateMatch(
	ctx context.Context, match_type string, blue_team, red_team []string,
) (db.Match, error) {
	if match_type != "indoor" && match_type != "beach" {
		return db.Match{}, fmt.Errorf("invalid match type, expected `indoor` or `beach`, got `%s`", match_type)
	}

	if len(blue_team) < 1 || len(red_team) < 1 {
		return db.Match{}, fmt.Errorf("each team must have at least one player")
	}

	match, err := s.queries.CreateMatch(ctx, db.CreateMatchParams{
		MatchType: match_type,
		Season: int32(time.Now().UTC().Year()),
	})
	if err != nil {
		return db.Match{}, fmt.Errorf("failed to create match: %w", err)
	}

	// Add players to match
	for _, player_name := range blue_team {
		// Check if player exists
		player, err := s.queries.GetPlayerByName(ctx, player_name)
		if err != nil {
			return db.Match{}, fmt.Errorf("player not found: %w", err)
		}
		if err := s.addPlayerToMatch(ctx, match, player, "blue"); err != nil {
			return db.Match{}, fmt.Errorf("failed to add players to blue team: %w", err)
		}
	}
	for _, player_name := range red_team {
		// Check if player exists
		player, err := s.queries.GetPlayerByName(ctx, player_name)
		if err != nil {
			return db.Match{}, fmt.Errorf("player not found: %w", err)
		}
		if err := s.addPlayerToMatch(ctx, match, player, "red"); err != nil {
			return db.Match{}, fmt.Errorf("failed to add players to red team: %w", err)
		}
	}

	return match, nil
}

func (s *MatchService) GetMatch(ctx context.Context, matchID uuid.UUID) (db.Match, error) {
	match, err := s.queries.GetMatch(ctx, matchID)
	if err != nil {
		return db.Match{}, fmt.Errorf("failed to retrieve match: %w", err)
	}

	return match, nil
}

func (s *MatchService) ListMatchesBySeason(
	ctx context.Context, match_type string, season int32,
) ([]db.Match, error) {
	matches, err := s.queries.ListMatchesBySeason(ctx, db.ListMatchesBySeasonParams{
		MatchType: match_type,
		Season: season,
	})
	if err != nil {
		return []db.Match{}, fmt.Errorf("failed to list seasonal matches: %w", err)
	}

	return matches, nil
}

func (s *MatchService) ListUncompletedMatches(ctx context.Context) ([]db.Match, error) {
	matches, err := s.queries.ListUncompletedMatches(ctx)
	if err != nil {
		return []db.Match{}, fmt.Errorf("failed to list uncompleted matches: %w", err)
	}

	return matches, nil
}

func (s *MatchService) DeleteUncompletedMatch(ctx context.Context, matchID uuid.UUID) error {
	match, err := s.queries.GetMatch(ctx, matchID)
	if err != nil {
		return fmt.Errorf("failed to retrieve match: %w", err)
	}

	if match.IsCompleted {
		return fmt.Errorf("cannot delete a match that was already registered")
	}

	if err := s.queries.DeleteUncompletedMatch(ctx, matchID); err != nil {
		return fmt.Errorf("failed to delete match: %w", err)
	}

	return nil
}

func (s *MatchService) GetMatchPlayers(ctx context.Context, matchID uuid.UUID) ([]db.GetMatchPlayersRow, error) {
	match_players, err := s.queries.GetMatchPlayers(ctx, matchID)
	if err != nil {
		return []db.GetMatchPlayersRow{}, fmt.Errorf("failed to retrieve players from match: %w", err)
	}

	return match_players, nil
}
