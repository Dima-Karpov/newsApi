package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"newsApi/internal/domain"
	"time"
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
func (r *NewsRepository) GetNews(page, pageSize int, fromDateStr, toDateStr, temp *string) ([]domain.NewsList, int, error) {
	var newsList []domain.NewsList

	// Парсим даты, если они переданы
	var fromDate, toDate *time.Time

	if fromDateStr != nil {
		parsedFromDate, err := time.Parse("2006-01-02", *fromDateStr)
		if err == nil {
			fromDate = &parsedFromDate
		}
	}

	if toDateStr != nil {
		parsedToDate, err := time.Parse("2006-01-02", *toDateStr)
		if err == nil {
			toDate = &parsedToDate
		}
	}

	// Формируем запрос
	query := r.db.Model(&domain.NewsList{})

	// Фильтрация по дате
	if fromDate != nil {
		query = query.Where("published_at >= ? ", *fromDate)
	}
	if toDate != nil {
		query = query.Where("published_at <= ? ", *toDate)
	}

	// Фильтрация по названию, если temp передан (поиск по подстроке, нечувствительный к регистру)
	if temp != nil {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+*temp+"%")
	}

	// Подсчитываем общее количество новостей после фильтрации
	var totalCount int64
	err := query.Count(&totalCount).Error
	if err != nil {
		return newsList, 0, errors.New("Error counting news: " + err.Error())
	}

	// Определяем порядок сортировки
	query = query.Order("published_at DESC")

	// Вычисляем смещение (offset) для пагинации
	err = query.Limit(pageSize).Offset((page - 1) * pageSize).Find(&newsList).Error
	if err != nil {
		return newsList, 0, errors.New("Error getting news: " + err.Error())
	}

	return newsList, int(totalCount), nil
}

// GetNew возвращает статью по id
func (r *NewsRepository) GetNew(id uuid.UUID) (domain.NewsList, error) {
	var news domain.NewsList

	err := r.db.Where("id = ?", id).First(&news).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return news, errors.New("News not found")
		}
		return news, errors.New("Error getting news: " + err.Error())
	}

	return news, nil
}
