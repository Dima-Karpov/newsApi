package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"newsApi/internal/domain"
)

type News interface {
	Save(news *domain.RSSItem) error
	GetNews(page, pageSize int, fromDateStr, toDateStr, temp *string) ([]domain.NewsList, int, error)
	GetNew(id uuid.UUID) (domain.NewsList, error)
}

type Repository struct {
	News
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		News: NewNews(db),
	}
}
