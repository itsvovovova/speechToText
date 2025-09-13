package main

import (
	"context"
	"fmt"
	"net/http"
	"speechToText/src/api"
	"speechToText/src/auth"
	"speechToText/src/config"
	"speechToText/src/consumer"
	"speechToText/src/metrics"
	"github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// @title Speech to Text API
// @version 1.0
// @description API для преобразования речи в текст
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	metrics := metrics.NewMetrics()
	var r = chi.NewRouter()
	ctx := context.Background()
	r.Use(metrics.Middleware)
	
	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", config.CurrentConfig.Server.Port)),
	))
	
	// API routes
	r.With(auth.Middleware).Get("/status", api.Status)
	r.With(auth.Middleware).Get("/result", api.Result)
	r.With(auth.Middleware).Post("/audio", api.Audio)
	r.With(auth.Middleware).Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Post("/login", api.Login)
	r.Post("/register", api.Register)
	go func() {
		if err := consumer.ReceiveMessage("queue", ctx); err != nil {
			fmt.Println("consumer error:", err)
		}
	}()
	err := http.ListenAndServe(":"+config.CurrentConfig.Server.Port, r)
	if err != nil {
		panic(err)
	}
}
