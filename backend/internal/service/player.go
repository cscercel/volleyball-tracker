package service

import (
	"context"
	"fmt"

	"github.com/cscercel/volleyball-tracker/internal/db"
	"github.com/google/uuid"
)

type PlayerService struct {
	queries	*db.Queries
}

type PlayerWithStats struct {
	PlayerName	string			`json:"player_name"`
	Stats		[]db.PlayerStat	`json:"stats"`
}

func NewPlayerService(queries *db.Queries) *PlayerService {
	return &PlayerService{queries: queries}
}

func (s *PlayerService) GetPlayerCareer(ctx context.Context, playerID uuid.UUID) (PlayerWithStats, error) {
	player, err := s.queries.GetPlayer(ctx, playerID)
	if err != nil {
		return PlayerWithStats{}, fmt.Errorf("unable to get player: %w", err)
	}
	
	// Get Stats
	stats, err := s.queries.GetPlayerStats(ctx, playerID)
	if err != nil {
		return PlayerWithStats{}, fmt.Errorf("could not load player stats: %w", err)
	}

	return PlayerWithStats{
		PlayerName: player.Name,
		Stats: stats,
	}, nil
}

func (s *PlayerService) GetPlayerSeason(
	ctx context.Context, playerID uuid.UUID, match_type string, season int,
) (PlayerWithStats, error) {
	player, err := s.queries.GetPlayer(ctx, playerID)
	if err != nil {
		return PlayerWithStats{}, fmt.Errorf("no player found: %w", err)
	}
	
	// Get Stats
	params := db.GetPlayerSeasonalStatsParams{
		PlayerID: playerID,
		MatchType: match_type,
		Season: int32(season),
	}

	stats, err := s.queries.GetPlayerSeasonalStats(ctx, params)
	if err != nil {
		return PlayerWithStats{}, fmt.Errorf("could not load player stats: %w", err)
	}

	return PlayerWithStats{
		PlayerName: player.Name,
		Stats: []db.PlayerStat{stats},
	}, nil
}

func (s *PlayerService) CreatePlayer(
	ctx context.Context, playerName, match_type string, season int,
) (PlayerWithStats, error) {
	player, err := s.queries.CreatePlayer(ctx, playerName)
	if err != nil {
		return PlayerWithStats{}, fmt.Errorf("could not create player: %w", err)		
	}

	// Initialize stats
	params := db.CreatePlayerStatsParams{
		PlayerID: player.ID,
		MatchType: match_type,
		Season: int32(season),
	}

	stats, err := s.queries.CreatePlayerStats(ctx, params)
	if err != nil {
		return PlayerWithStats{}, fmt.Errorf("could not initialize player stats: %w", err)
	}

	return PlayerWithStats{
		PlayerName: player.Name,
		Stats: []db.PlayerStat{stats},
	}, nil
}

func (s *PlayerService) EditPlayerName(ctx context.Context, playerID uuid.UUID, name string) (db.Player, error) {
	params := db.EditPlayerNameParams{
		ID: playerID,
		Name: name,
	}

	player, err := s.queries.EditPlayerName(ctx, params)
	if err != nil {
		return db.Player{}, fmt.Errorf("could not change player name: %w", err)
	}

	return player, nil
}

func (s *PlayerService) DeletePlayer(ctx context.Context, playerID uuid.UUID) error {
	if err := s.queries.DeletePlayer(ctx, playerID); err != nil {
		return fmt.Errorf("could not delete player: %w", err)
	}
	
	return nil
}

func (s *PlayerService) ListRoster(ctx context.Context) ([]db.Player, error) {
	roster, err := s.queries.ListPlayers(ctx)
	if err != nil {
		return []db.Player{}, fmt.Errorf("could not list roster: %w", err)
	}

	// No need to fetch stats for entire roster, will be ONLY season specific
	// We just want the names here for a dropdown list on the client side
	return roster, nil
}

func (s *PlayerService) ListSeasonalRoster(
	ctx context.Context, match_type string, season int,
) ([]PlayerWithStats, error) {
	roster, err := s.queries.ListPlayers(ctx)
	if err != nil {
		return []PlayerWithStats{}, fmt.Errorf("could not list roster: %w", err)
	}

	params := db.ListSeasonalStatsParams{
		MatchType: match_type,
		Season: int32(season),
	}
	
	// Merge stats
	stats, err := s.queries.ListSeasonalStats(ctx, params)
	if err != nil {
		return []PlayerWithStats{}, fmt.Errorf("could not load seasonal stats: %w", err)
	}

	for _, player := range roster {
		
	}

	return stats, nil
}
