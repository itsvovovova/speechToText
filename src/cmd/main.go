package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	var ctx = context.Background()
	fmt.Println(ctx)
	var r = chi.NewRouter()
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}
}
