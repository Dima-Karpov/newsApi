package repository

import (
	"gorm.io/gorm"
	"newsApi/internal/domain"
)

type News interface {
	Save(news *domain.RSSItem) error
	GetNews(page, pageSize int) ([]domain.NewsList, int, error)
}

type Repository struct {
	News
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		News: NewNews(db),
	}
}
