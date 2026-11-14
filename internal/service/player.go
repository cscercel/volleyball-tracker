package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cscercel/volleyball-tracker/internal/db"
)

type PlayerService struct {
	queries *db.Queries
}

func NewPlayerService(queries *db.Queries) *PlayerService {
	return &PlayerService{queries: queries}
}

func (s *PlayerService) CreatePlayer(ctx context.Context, name string) (db.Player, error) {
	player, err := s.queries.CreatePlayer(ctx, name)
	if err != nil {
		return db.Player{}, fmt.Errorf("failed to create player: %w", err)
	}

	// Initialize player stats for both indoor and beach this season
	for _, match_type := range []string{"indoor", "beach"} {
		_, err := s.queries.UpsertPlayerStats(ctx, db.UpsertPlayerStatsParams{
			PlayerID:  player.ID,
			MatchType: match_type,
			Season:    int64(time.Now().UTC().Year()),
		})
		if err != nil {
			return db.Player{}, fmt.Errorf("failed to create stats for %v:%w", match_type, err)
		}
	}

	return player, nil
}

func (s *PlayerService) ListPlayers(ctx context.Context) ([]db.Player, error) {
	players, err := s.queries.ListPlayers(ctx)
	if err != nil {
		return []db.Player{}, fmt.Errorf("failed to load players: %w", err)
	}

	return players, nil
}

func (s *PlayerService) UpdatePlayerName(
	ctx context.Context, playerID int64, newName string,
) (db.Player, error) {
	player, err := s.queries.UpdatePlayerName(ctx, db.UpdatePlayerNameParams{
		ID:   playerID,
		Name: newName,
	})
	if err != nil {
		return db.Player{}, fmt.Errorf("failed to update player name: %w", err)
	}

	return player, nil
}

func (s *PlayerService) DeletePlayer(ctx context.Context, playerID int64) error {
	if err := s.queries.DeletePlayer(ctx, playerID); err != nil {
		return fmt.Errorf("failed to delete player: %w", err)
	}

	return nil
}

func (s *PlayerService) GetPlayerByID(
	ctx context.Context, playerID int64, match_type string, season int64,
) (db.GetPlayerStatsByIDRow, error) {
	player, err := s.queries.GetPlayerStatsByID(ctx, db.GetPlayerStatsByIDParams{
		PlayerID:  playerID,
		MatchType: match_type,
		Season:    season,
	})
	if err != nil {
		return db.GetPlayerStatsByIDRow{}, fmt.Errorf("failed to get player: %w", err)
	}

	return player, nil
}

func (s *PlayerService) GetPlayerByName(
	ctx context.Context, playerName, match_type string, season int64,
) (db.GetPlayerStatsByNameRow, error) {
	player, err := s.queries.GetPlayerStatsByName(ctx, db.GetPlayerStatsByNameParams{
		Name:      playerName,
		MatchType: match_type,
		Season:    season,
	})
	if err != nil {
		return db.GetPlayerStatsByNameRow{}, fmt.Errorf("failed to get player: %w", err)
	}

	return player, nil
}

func (s *PlayerService) GetLeaderboard(
	ctx context.Context, match_type string, season int64,
) ([]db.GetLeaderboardRow, error) {
	leaderboard, err := s.queries.GetLeaderboard(ctx, db.GetLeaderboardParams{
		MatchType: match_type,
		Season:    season,
	})
	if err != nil {
		return []db.GetLeaderboardRow{}, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	return leaderboard, nil
}

// From matches.sql -> felt it was better here
func (s *PlayerService) GetPlayerSeasonalMatches(
	ctx context.Context, playerID int64, match_type string, season, page, pageLimit int64,
) ([]db.GetPlayerSeasonalMatchesRow, bool, error) {
	if page < 1 {
		page = 1
	}
	if pageLimit < 1 {
		pageLimit = 5
	}
	if pageLimit > 50 {
		pageLimit = 50 // Cap
	}

	limit := pageLimit + 1 // fetch one extra to detect next page
	offset := (page - 1) * pageLimit

	player_matches, err := s.queries.GetPlayerSeasonalMatches(ctx, db.GetPlayerSeasonalMatchesParams{
		PlayerID:  playerID,
		MatchType: match_type,
		Season:    season,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return []db.GetPlayerSeasonalMatchesRow{}, false, fmt.Errorf("failed to retrieve player matches: %w", err)
	}

	hasMore := int64(len(player_matches)) > pageLimit
	if hasMore {
		player_matches = player_matches[:pageLimit]
	}

	return player_matches, hasMore, nil
}
