package usecase

import (
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
	"newsApi/internal/domain"
	"newsApi/internal/repository"
	"time"
)

type RSSHandler struct {
	repo         *repository.Repository
	rssURLs      []string
	pollInterval time.Duration
	ctx          context.Context
}

func NewRSSHandler(
	repo *repository.Repository,
	urls []string,
	pollInterval time.Duration,
	ctx context.Context,
) *RSSHandler {
	return &RSSHandler{
		repo:         repo,
		rssURLs:      urls,
		pollInterval: pollInterval,
		ctx:          ctx,
	}
}

func (h *RSSHandler) Start() {
	parser := gofeed.NewParser()
	ticker := time.NewTicker(h.pollInterval)
	defer ticker.Stop()

	// Выполняем первый запрос сразу
	h.processFeeds(parser)

	for {
		select {
		case <-h.ctx.Done():
			fmt.Println("RSS handler shutting down...")
			return
		case <-ticker.C:
			h.processFeeds(parser)
		}
	}
}

func (h *RSSHandler) processFeeds(parser *gofeed.Parser) {
	for _, url := range h.rssURLs {
		fmt.Printf("Fetching RSS feed: %s\n", url)

		// Парсим RSS-ленту
		parsedFeed, err := parser.ParseURL(url)
		if err != nil {
			fmt.Printf("Error parsing RSS feed: %s\n", err)
			continue
		}

		// Обрабатываем записи из RSS-ленты
		for _, item := range parsedFeed.Items {
			rssItem := &domain.RSSItem{
				Title:       item.Title,
				Description: item.Description,
				PublishedAt: *item.PublishedParsed,
				Link:        item.Link,
			}

			// Сохраняем новость в базу через репозиторий
			if err := h.repo.Save(rssItem); err != nil {
				fmt.Printf("Error saving RSS feed: %s\n", err)
				continue
			}
		}
	}
}
