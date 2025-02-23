package service

import (
	"github.com/google/uuid"
	"newsApi/internal/domain"
	"newsApi/internal/repository"
)

type NewsList interface {
	GetNews(page, pageSize int, fromDateStr, toDateStr, temp *string) ([]domain.NewsList, int, error)
	GetNew(id uuid.UUID) (domain.NewsList, error)
}

type Service struct {
	NewsList
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		NewsList: NewNewsService(repos),
	}
}
