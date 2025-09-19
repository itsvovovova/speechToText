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
	ctx := context.Background()
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
	go func() {
		if err := consumer.ReceiveMessage("queue", ctx); err != nil {
			fmt.Println("consumer error:", err)
		}
	}()
	server := &http.Server{
		Addr:    ":" + config.CurrentConfig.Server.Port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		fmt.Printf("Server started on port %s\n", config.CurrentConfig.Server.Port)
		fmt.Printf("Swagger UI: http://localhost:%s/swagger/index.html\n", config.CurrentConfig.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	<-quit
	fmt.Println("Received shutdown signal, stopping server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Error stopping server: %v\n", err)
	} else {
		fmt.Println("Server stopped successfully")
	}
}
