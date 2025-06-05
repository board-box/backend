package game

import (
	"context"

	pgx "github.com/jackc/pgx/v5"
)

type Service struct {
	repo *repository
}

func NewService(db *pgx.Conn) *Service {
	return &Service{
		repo: newRepository(db),
	}
}

func (s *Service) ListGames(ctx context.Context) ([]Game, error) {
	return s.repo.listGames(ctx)
}

func (s *Service) GetGame(ctx context.Context, id int64) (Game, error) {
	return s.repo.getGameById(ctx, id)
}

func (s *Service) GetGames(ctx context.Context, ids []int64) ([]Game, error) {
	return s.repo.getGamesByIds(ctx, ids)
}

func (s *Service) CreateGame(ctx context.Context, game Game) (int64, error) {
	return s.repo.createGame(ctx, game)
}

func (s *Service) UpdateGame(ctx context.Context, game Game) error {
	return s.repo.updateGame(ctx, game)
}

func (s *Service) DeleteGame(ctx context.Context, id int64) error {
	return s.repo.deleteGame(ctx, id)
}
