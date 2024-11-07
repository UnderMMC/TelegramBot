package service

import (
	"context"
	"log"
)

type BotRepository interface {
	SearchAllByFilter(ctx context.Context) (string, error)
}

type BotService struct {
	repository BotRepository
}

func NewBotService(repository BotRepository) *BotService {
	return &BotService{repository: repository}
}

func (s *BotService) ShowAllRepairment(ctx context.Context) (string, error) {
	result, err := s.repository.SearchAllByFilter(ctx)
	if err != nil {
		log.Printf("Ошибка обращения к слою репозитория: %v", err)
		return "", err
	}

	return result, nil
}
