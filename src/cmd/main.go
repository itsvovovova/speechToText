package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"speechToText/src/api"
	"speechToText/src/auth"
	"speechToText/src/cache"
	"speechToText/src/config"
	"speechToText/src/consumer"
	"speechToText/src/db"
	"speechToText/src/docs"
	appmetrics "speechToText/src/metrics"
	"speechToText/src/pkg/closer"
	"sync"
	"syscall"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// @title Speech to Text API
// @version 1.0
// @description Speech to Text API
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cl := closer.New()

	sqlDB, err := db.NewDB(config.CurrentConfig.Database)
	if err != nil {
		log.Fatalf("database init: %v", err)
	}
	cl.Add(sqlDB.Close)
	store := db.NewStore(sqlDB)

	sessionProvider := cache.NewRedisSessionProvider(config.CurrentConfig.Redis.Host)
	cl.Add(sessionProvider.Close)
	sessionManager := cache.NewRedisSessionManager("session_id", sessionProvider, int64(math.Pow10(5)))

	producer := consumer.NewProducer(config.CurrentConfig.RabbitMQ.Url)
	cl.Add(producer.Close)

	cons := consumer.NewConsumer(store)

	handlers := api.NewHandlers(store, sessionManager, producer)
	authMiddleware := auth.NewMiddleware(sessionManager)

	docs.SwaggerInfo.Title = "Speech to Text API"
	docs.SwaggerInfo.Description = "Speech to Text API with pagination"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + config.CurrentConfig.Server.Port
	docs.SwaggerInfo.BasePath = "/"

	m := appmetrics.NewMetrics()
	r := chi.NewRouter()
	r.Use(m.Middleware)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", config.CurrentConfig.Server.Port)),
	))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Post("/register", handlers.Register)
	r.Post("/login", handlers.Login)

	r.With(authMiddleware).Post("/logout", handlers.Logout)
	r.With(authMiddleware).Post("/audio", handlers.Audio)
	r.With(authMiddleware).Get("/status", handlers.Status)
	r.With(authMiddleware).Get("/result", handlers.Result)
	r.With(authMiddleware).Get("/tasks", handlers.Tasks)
	r.With(authMiddleware).Delete("/tasks/{id}", handlers.DeleteTask)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := cons.Receive("queue", ctx); err != nil {
			log.Println("consumer error:", err)
		}
	}()

	server := &http.Server{
		Addr:         ":" + config.CurrentConfig.Server.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Server started on port %s\n", config.CurrentConfig.Server.Port)
		log.Printf("Swagger UI: http://localhost:%s/swagger/index.html\n", config.CurrentConfig.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Received shutdown signal, stopping server...")

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error stopping server: %v\n", err)
	}

	wg.Wait()
	cl.Close()
	log.Println("Server stopped successfully")
}
