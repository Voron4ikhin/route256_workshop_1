package main

import (
	"log"
	"route256/cart/internal/pkg/app/http"
)

func main() {
	app := http.NewApp()
	log.Fatal(app.Run())
}
