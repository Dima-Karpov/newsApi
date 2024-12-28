package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"newsApi/internal/domain"
)

type NewsRepository struct {
	db *gorm.DB
}

func NewNews(db *gorm.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

// Save сохраняет новость в базу
func (r *NewsRepository) Save(news *domain.RSSItem) error {
	// Проверяем уникальность по ссылке
	existing := domain.NewsList{}
	err := r.db.Where("link = ?", news.Link).First(&existing).Error

	// Если ошибка записи не найдена (gorm.ErrRecordNotFound), значит записи нет
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("error checking if news exists: " + err.Error())
	}

	// Если запись не найдена, то создаем новую
	if err == gorm.ErrRecordNotFound {
		result := &domain.NewsList{
			ID:          uuid.New(),
			Title:       news.Title,
			Description: news.Description,
			Link:        news.Link,
			PublishedAt: news.PublishedAt,
		}

		// Сохраняем новость в базу данных
		if err := r.db.Create(&result).Error; err != nil {
			return errors.New("error saving news: " + err.Error())
		}
	} else {
		// Если запись существует, мы можем вернуть ошибку или просто ничего не делать
		return errors.New("news already exists with link: " + news.Link)
	}

	return nil
}

// GetNews возвращает все стати
func (r *NewsRepository) GetNews(page, pageSize int) ([]domain.NewsList, int, error) {
	var newsList []domain.NewsList

	// Вычисляем смещение (offset) для пагинации
	offset := (page - 1) * pageSize

	// Получаем новости с ограничением (limit) и смещением (offset)
	err := r.db.Limit(pageSize).Offset(offset).Find(&newsList).Error
	if err != nil {
		return newsList, 0, errors.New("Error getting news: " + err.Error())
	}

	// Подсчитываем общее количество новостей
	var totalCount int64
	err = r.db.Model(&domain.NewsList{}).Count(&totalCount).Error
	if err != nil {
		return newsList, 0, errors.New("Error counting news: " + err.Error())
	}

	return newsList, int(totalCount), nil
}
