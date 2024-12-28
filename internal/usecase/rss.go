package usecase

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"newsApi/internal/domain"
	"newsApi/internal/repository"
	"time"
)

// Логика работы с RSS

type RSSFeed struct {
	URL           string
	PollInterval  time.Duration
	ProcessedNews map[string]bool // Для предотвращения дублирования
}

// Структура, отвечающая за управление несколькими RSS-лентами.
type RSSHandler struct {
	Feeds  []RSSFeed
	Output chan *domain.RSSItem
	repo   *repository.Repository
}

// NewRSSHandler обработчик RSS
func NewRSSHandler(repo *repository.Repository, urls []string, pollInterval time.Duration) *RSSHandler {
	feeds := make([]RSSFeed, len(urls))
	for i, url := range urls {
		feeds[i] = RSSFeed{
			URL:           url,
			PollInterval:  pollInterval,
			ProcessedNews: make(map[string]bool),
		}
	}

	return &RSSHandler{
		Feeds:  feeds,
		Output: make(chan *domain.RSSItem),
		repo:   repo,
	}
}

// Start запускает обработку всех RSS-лент
func (h *RSSHandler) Start() {
	for _, feed := range h.Feeds {
		go h.processFeed(feed)
	}
}

// processFeed выполняет периодическое чтение одной RSS-ленты
func (h *RSSHandler) processFeed(feed RSSFeed) {
	parser := gofeed.NewParser()

	for {
		fmt.Printf("Fetching RSS feed: %s\n", feed.URL)

		// Парсим RSS-ленту
		parsedFeed, err := parser.ParseURL(feed.URL)
		if err != nil {
			fmt.Printf("Error parsing RSS feed: %s\n", err)
			time.Sleep(feed.PollInterval)
			continue
		}

		// Обрабатываем записи из RSS-ленты
		for _, item := range parsedFeed.Items {
			if _, exists := feed.ProcessedNews[item.GUID]; exists {
				continue // Пропускаем уже обработанные публикации
			}

			// Создаем новость
			rssItem := &domain.RSSItem{
				Title:       item.Title,
				Description: item.Description,
				PublishedAt: *item.PublishedParsed,
				Link:        item.Link,
			}

			// Отправляем новость в канал
			h.Output <- rssItem

			// Помечаем новость как обработанную
			feed.ProcessedNews[item.GUID] = true

			// Сохраняем новость в базу через репозиторий
			if err := h.repo.Save(rssItem); err != nil {
				fmt.Printf("Error saving RSS feed: %s\n", err)
				continue
			}
		}

		// Ждем следующий цикл опроса
		time.Sleep(feed.PollInterval)
	}
}
