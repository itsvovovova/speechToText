package main

import (
	"github.com/go-chi/chi/v5"
	"math"
	"net/http"
	"speechToText/src/cache"
)

var SessionManager *cache.RedisSessionManager

func main() {
	var sessionProvider, err = cache.NewRedisSessionProvider(":8080")
	if err != nil {
		panic(err)
	}
	SessionManager = cache.NewRedisSessionManager("session_id", sessionProvider, int64(math.Pow10(5)))
	var r = chi.NewRouter()
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}
}
