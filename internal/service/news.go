package service

import (
	"github.com/google/uuid"
	"newsApi/internal/domain"
	"newsApi/internal/repository"
)

type NewsService struct {
	repo repository.News
}

func NewNewsService(repo repository.News) *NewsService {
	return &NewsService{repo: repo}
}

func (s *NewsService) GetNews(page, pageSize int, fromDateStr, toDateStr, temp *string) ([]domain.NewsList, int, error) {
	return s.repo.GetNews(page, pageSize, fromDateStr, toDateStr, temp)
}
func (s *NewsService) GetNew(id uuid.UUID) (domain.NewsList, error) {
	return s.repo.GetNew(id)
}
