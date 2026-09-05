package main

import (
	"log"
	"net/http"
	"route256/loms/internal/pkg/handlers"
	"route256/loms/internal/pkg/repository"
	"route256/loms/internal/pkg/services"
)

func main() {
	repo := repository.NewDumbRepo()
	stocksHandler := handlers.NewStocksHandler(services.NewStocksService(repo))
	http.HandleFunc("/stocks", stocksHandler.Handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
