package service

import (
	"newsApi/internal/domain"
	"newsApi/internal/repository"
)

type NewsList interface {
	GetNews(page, pageSize int) ([]domain.NewsList, int, error)
}

type Service struct {
	NewsList
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		NewsList: NewNewsService(repos),
	}
}
