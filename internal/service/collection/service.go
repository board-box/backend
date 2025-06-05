package collection

import (
	"context"
	"errors"

	"github.com/board-box/backend/internal/service/game"
	pgx "github.com/jackc/pgx/v5"
)

var (
	ErrCollectionNotFound = errors.New("collection not found")
	ErrForbidden          = errors.New("forbidden")
)

type Service struct {
	repo    *repository
	gameSvc *game.Service
}

func NewService(db *pgx.Conn, gameSvc *game.Service) *Service {
	return &Service{
		repo:    newRepository(db),
		gameSvc: gameSvc,
	}
}

func (s *Service) ListCollections(ctx context.Context, userID int64) ([]Collection, error) {
	return s.repo.listCollections(ctx, userID)
}

func (s *Service) GetCollection(ctx context.Context, collectionID, userID int64) (Collection, error) {
	return s.repo.getCollection(ctx, collectionID, userID)
}

func (s *Service) CreateCollection(ctx context.Context, userID int64, req Collection) (Collection, error) {
	if req.Name == "" {
		return Collection{}, errors.New("collection name cannot be empty")
	}
	return s.repo.createCollection(ctx, userID, req)
}

func (s *Service) UpdateCollection(ctx context.Context, collectionID, userID int64, req Collection) (Collection, error) {
	if req.Name == "" {
		return Collection{}, errors.New("collection name cannot be empty")
	}
	return s.repo.updateCollection(ctx, collectionID, userID, req)
}

func (s *Service) DeleteCollection(ctx context.Context, collectionID, userID int64) error {
	return s.repo.deleteCollection(ctx, collectionID, userID)
}

func (s *Service) AddGameToCollection(ctx context.Context, collectionID, gameID, userID int64) error {
	_, err := s.gameSvc.GetGame(ctx, gameID)
	if err != nil {
		return err
	}

	return s.repo.addGameToCollection(ctx, collectionID, gameID, userID)
}

func (s *Service) RemoveGameFromCollection(ctx context.Context, collectionID, gameID, userID int64) error {
	return s.repo.removeGameFromCollection(ctx, collectionID, gameID, userID)
}
