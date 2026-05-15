package service

import (
	"context"
	"errors"

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

func (s *PlayerService) GetPlayerCareer(ctx context.Context, playerID uuid.UUID) (PlayerWithStats, error) {
	player, err := s.queries.GetPlayer(ctx, playerID)
	if err != nil {
		return PlayerWithStats{}, errors.New("no player found")
	}
	
	// Get Stats
	stats, err := s.queries.GetPlayerStats(ctx, playerID)
	if err != nil {
		return PlayerWithStats{}, errors.New("could not load player stats")
	}

	return PlayerWithStats{
		PlayerName: player.Name,
		Stats: stats,
	}, nil
}

func (s *PlayerService) GetPlayerSeason(
	ctx context.Context, playerID uuid.UUID, match_type string, season int,
) (PlayerWithStats, error) {
	// TODO
}
