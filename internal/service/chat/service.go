package chat

import (
	"context"

	"github.com/samber/lo"
)

type Service struct {
	client  *client
	history map[int64][]message
}

func NewService(apiKey string) *Service {
	return &Service{
		client:  newClient(apiKey),
		history: make(map[int64][]message),
	}
}

func (s *Service) Chat(_ context.Context, userID int64, msg string) ([]string, error) {
	if _, ok := s.history[userID]; !ok {
		s.history[userID] = []message{
			{Role: "system", Content: "Ты эксперт по настольным играм и сможешь подобрать нужную игру по запросу. Отвечай коротко и по делу без схем и списков."},
		}
	}
	s.history[userID] = append(s.history[userID], message{Role: "user", Content: msg})

	answer, err := s.client.chat(s.history[userID])
	if err != nil {
		return nil, err
	}
	s.history[userID] = append(s.history[userID], answer)

	return lo.Map(s.history[userID][1:], func(msg message, index int) string {
		return msg.Content
	}), nil
}
