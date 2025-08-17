package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"speechToText/src/api"
	"speechToText/src/auth"
	"speechToText/src/cache"
	"speechToText/src/config"
	"speechToText/src/consumer"
)

var SessionManager *cache.RedisSessionManager

func main() {
	var r = chi.NewRouter()
	ctx := context.Background()
	r.Get("/status", api.Status)
	r.Get("/result", api.Result)
	r.Post("/audio", api.Audio)
	r.Use(auth.Middleware)
	err := http.ListenAndServe(config.CurrentConfig.Server.Port, r)
	if err = consumer.ReceiveMessage("queue", ctx); err != nil {
		fmt.Println("consumer error: ", err)
	}
	if err != nil {
		panic(err)
	}
}
