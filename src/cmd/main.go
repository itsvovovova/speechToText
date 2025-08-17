package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"speechToText/src/api"
	"speechToText/src/auth"
	"speechToText/src/cache"
)

var SessionManager *cache.RedisSessionManager

func main() {
	var r = chi.NewRouter()
	r.Get("/status", api.Status)
	r.Get("/result", api.Result)
	r.Post("/audio", api.Audio)
	r.Use(auth.Middleware)
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}
}
