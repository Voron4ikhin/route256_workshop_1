package main

import (
	"log"
	"route256/loms/internal/pkg/app/http"
	"route256/loms/internal/pkg/config"
)

func main() {
	cfg := config.NewFromFlags()
	app := http.NewApp(cfg)
	log.Fatal(app.Run())
}
