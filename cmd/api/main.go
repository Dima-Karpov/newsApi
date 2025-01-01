package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"newsApi/configs"
	"newsApi/internal/delivery/router"
	delivery "newsApi/internal/delivery/server"
	"newsApi/internal/repository"
	"newsApi/internal/service"
	"newsApi/internal/usecase"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// @title News API
// @version 1.0
// @description API Server for NewsApi Application

// @host localhost:8000
// @BasePath /

func main() {
	cfg, err := configs.LoadConfig("configs/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %s\n", err.Error())
	}

	// Загрузка переменных окружения из файла .env
	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	rssURLs := cfg.RSS
	requestPeriod := time.Duration(cfg.RequestPeriod) * time.Minute
	cfgDB := repository.Config{DSN: dsn}

	db, err := repository.NewPostgresDB(cfgDB)
	if err != nil {
		fmt.Printf("Failed to initialize db: %s\n", err.Error())
		return
	}

	// Создаем зависимости
	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlerRouterNews := router.NewHandler(services)

	// Создаём HTTP сервер
	srv := new(delivery.Server)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Создаём RSS обработчик с поддержкой контекста
	handlerRSS := usecase.NewRSSHandler(repo, rssURLs, requestPeriod, ctx)

	var wg sync.WaitGroup

	// Запускаем HTTP сервер
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Run("8000", handlerRouterNews.InitRouter()); err != nil {
			fmt.Printf("Error occurred while running HTTP server: %s\n", err.Error())
		}
	}()

	// Запускаем RSS обработчик
	wg.Add(1)
	go func() {
		defer wg.Done()
		handlerRSS.Start()
	}()

	// Ожидаем сигнала завершения
	<-ctx.Done()
	fmt.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Error shutting down server: %s\n", err.Error())
	}

	// Ожидаем завершения всех горутин
	wg.Wait()
	fmt.Println("Application stopped")
}
