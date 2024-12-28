package main

import (
	"fmt"
	"newsApi/internal/delivery/router"
	delivery "newsApi/internal/delivery/server"
	"newsApi/internal/repository"
	"newsApi/internal/service"
	"newsApi/internal/usecase"
	"time"
)

func main() {
	rssURLs := []string{
		"https://habr.com/ru/rss/hub/go/all/?fl=ru",
		"https://habr.com/ru/rss/best/daily/?fl=ru",
		"https://cprss.s3.amazonaws.com/golangweekly.com.xml",
	}
	cfgDB := repository.Config{
		DSN: "host=localhost port=5440 user=news password=OvoIpFrIL2VS dbname=api sslmode=disable",
	}

	db, err := repository.NewPostgresDB(cfgDB)
	if err != nil {
		fmt.Printf("Failed to initialize db: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	handlerRSS := usecase.NewRSSHandler(repo, rssURLs, 5*time.Minute)
	services := service.NewService(repo)
	handlerRouterNews := router.NewHandler(services)

	srv := new(delivery.Server)
	go func() {
		if err := srv.Run("8000", handlerRouterNews.InitRouter()); err != nil {
			fmt.Errorf("Error occured while running http sever: %s", err.Error())
		}
	}()

	go handlerRSS.Start()

	for item := range handlerRSS.Output {
		fmt.Printf("New item: %s\n", item.Title)
	}
}
