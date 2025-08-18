package main

import (
	"context"
	"fmt"
	"net/http"
	"speechToText/src/api"
	"speechToText/src/auth"
	"speechToText/src/config"
	"speechToText/src/consumer"

	"github.com/go-chi/chi/v5"
)

func main() {
	var r = chi.NewRouter()
	ctx := context.Background()
	r.With(auth.Middleware).Get("/status", api.Status)
	r.With(auth.Middleware).Get("/result", api.Result)
	r.With(auth.Middleware).Post("/audio", api.Audio)
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
