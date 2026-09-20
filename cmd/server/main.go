package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paraumir/shortener/internal/config"
	"github.com/paraumir/shortener/internal/handler"
	"github.com/paraumir/shortener/internal/service"
	"github.com/paraumir/shortener/internal/storage"
)

func main() {
	cfg := config.Load()

	var store storage.Storage
	var closeStorage func()

	switch cfg.StorageType {
	case "memory":
		store = storage.NewMemoryStorage()
		closeStorage = func() {}

	case "postgres":
		if cfg.DatabaseURL == "" {
			log.Fatal("DATABASE_URL is required for postgres storage")
		}

		ctx := context.Background()

		db, err := storage.NewPostgresPool(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("failed to connect to PostgreSQL: %v", err)
		}

		store = storage.NewPostgresStorage(db)
		closeStorage = db.Close

	default:
		log.Fatalf("unknown storage type: %s", cfg.StorageType)
	}

	defer closeStorage()

	svc := service.New(store)
	h := handler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", h.Shorten)
	mux.HandleFunc("/", h.GetOriginalURL)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("server started on :8080")
		log.Printf("storage: %s", cfg.StorageType)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}
