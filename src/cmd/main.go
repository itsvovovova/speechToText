package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"speechToText/src/cache"
)

var SessionManager *cache.RedisSessionManager

func main() {
	var r = chi.NewRouter()
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}
}
