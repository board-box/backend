package user

import (
	"context"
	"errors"
	"math/rand"

	"github.com/board-box/backend/internal/auth"
	pgx "github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrUserNotFound = errors.New("user not found")
)

type Service struct {
	repo *repository
	jwt  *auth.JWTManager
}

func NewService(db *pgx.Conn, jwt *auth.JWTManager) *Service {
	return &Service{
		repo: newRepository(db),
		jwt:  jwt,
	}
}

func (s *Service) Register(ctx context.Context, username, email, password string) error {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := User{ID: rand.Int63(), Username: username, Email: email, PasswordHash: string(hashed)} // nolint:gosec

	err := s.repo.saveUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.getUserByEmail(ctx, email)
	if err != nil {
		return "", ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrUnauthorized
	}

	token, err := s.jwt.GenerateToken(s.jwt.SecretKey, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Info(ctx context.Context, userID int64) (User, error) {
	user, err := s.repo.getUserByID(ctx, userID)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	return user, nil
}
