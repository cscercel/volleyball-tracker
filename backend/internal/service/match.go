package service

import (
	"context"
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

func (s *MatchService) GetMatch(ctx context.Context, match_id uuid.UUID) (db.Match, error) {
	match, err := s.queries.GetMatch(ctx, match_id)
	if err != nil {
		return db.Match{}, fmt.Errorf("unable to get match: %w", err)
	}

	return match, nil
}

func (s *MatchService) GetRegisteredMatches(ctx context.Context) ([]db.Match, error) {
	matches, err := s.queries.GetRegisteredMatches(ctx)
	if err != nil {
		return []db.Match{}, fmt.Errorf("unable to get registered matches: %w", err)
	}

	return matches, nil
}

func (s *MatchService) GetDrafts(ctx context.Context) ([]db.Match, error) {
	drafts, err := s.queries.GetDrafts(ctx)
	if err != nil {
		return []db.Match{}, fmt.Errorf("unable to get drafts: %w", err)
	}

	return drafts, nil
}

func (s *MatchService) GetSeasonMatches(ctx context.Context, match_type string, season int) ([]db.Match, error) {
	params := db.GetSeasonMatchesParams{
		MatchType: match_type,
		Season: int32(season),
	}

	matches, err := s.queries.GetSeasonMatches(ctx, params)
	if err != nil {
		return []db.Match{}, fmt.Errorf("unable to get seasonal matches: %w", err)
	}

	return matches, nil
}

func (s *MatchService) CreateMatch(
	ctx context.Context, match_type string, blue_team, red_team []db.Player
) (db.Match, error) {
	if match_type != "indoor" && match_type != "beach" {
		return db.Match{}, fmt.Errorf("invalid match_type: expected `indoor` or `beach`, got %s", match_type)
	}

	season := time.Now().UTC().Year()
	match, err := s.queries.CreateMatch(ctx, match_type, season)
	if err != nil {
		return db.Match{}, fmt.Errorf("could not create match: %w", err)
	}

	// Add players from blue_team
	for _, player := range blue_team {
		color := "blue"
		blue_team, err := s.queries.AddPlayerToMatch(ctx, match.ID, player.ID, color)
		if err != nil {
			return db.Match{}, fmt.Errorf("failed to add players from blue_team: %w", err)
		}
	}

	// idem. for red_team
	for _, player := range red_team {
		color := "red"
		red_team, err := s.queries.AddPlayerToMatch(ctx, match.ID, player.ID, color)
		if err != nil {
			return db.Match{}, fmt.Errorf("failed to add players from red_team: %w", err)
		}
	}

}
