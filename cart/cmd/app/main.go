package main

import (
	"log"
	"route256/cart/internal/pkg/app/http"
	"route256/cart/internal/pkg/config"
)

func main() {
	cfg := config.NewFromFlags()
	app := http.NewApp(cfg)
	log.Fatal(app.Run())
}
