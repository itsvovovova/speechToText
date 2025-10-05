package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"speechToText/src/api"
	"speechToText/src/auth"
	"speechToText/src/config"
	"speechToText/src/consumer"
	"speechToText/src/docs"
	"speechToText/src/metrics"
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
	docs.SwaggerInfo.Title = "Speech to Text API"
	docs.SwaggerInfo.Description = "Speech to Text API with pagination"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + config.CurrentConfig.Server.Port
	docs.SwaggerInfo.BasePath = "/"

	metrics := metrics.NewMetrics()
	var r = chi.NewRouter()
	r.Use(metrics.Middleware)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", config.CurrentConfig.Server.Port)),
	))
	r.With(auth.Middleware).Get("/status", api.Status)
	r.With(auth.Middleware).Get("/result", api.Result)
	r.With(auth.Middleware).Post("/audio", api.Audio)
	r.With(auth.Middleware).Get("/tasks", api.Tasks)
	r.With(auth.Middleware).Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Post("/login", api.Login)
	r.Post("/register", api.Register)
	// Graceful shutdown setup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := consumer.ReceiveMessage("queue", ctx); err != nil {
			fmt.Println("consumer error:", err)
		}
	}()

	// Start HTTP server
	server := &http.Server{
		Addr:    ":" + config.CurrentConfig.Server.Port,
		Handler: r,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("Server started on port %s\n", config.CurrentConfig.Server.Port)
		fmt.Printf("Swagger UI: http://localhost:%s/swagger/index.html\n", config.CurrentConfig.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Received shutdown signal, stopping server...")

	cancel() // Stop consumer
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Error stopping server: %v\n", err)
	}

	wg.Wait()
	fmt.Println("Server stopped successfully")
}
