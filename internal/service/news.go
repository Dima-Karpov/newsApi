package service

import (
	"newsApi/internal/domain"
	"newsApi/internal/repository"
)

type NewsService struct {
	repo repository.News
}

func NewNewsService(repo repository.News) *NewsService {
	return &NewsService{repo: repo}
}

func (s *NewsService) GetNews(page, pageSize int) ([]domain.NewsList, int, error) {
	return s.repo.GetNews(page, pageSize)
}
