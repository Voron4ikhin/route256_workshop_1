package http

import (
	"flag"
	"net/http"
	"os"
	"route256/loms/internal/pkg/handlers"
	"route256/loms/internal/pkg/repository"
	"route256/loms/internal/pkg/services"
)

type config struct {
	addr string
}

// envOrDefault returns the value of the env var if set, otherwise def.
// Flags below use this as their default, so the precedence is: flag > env > hardcoded default.
func envOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func newConfigFromFlags() config {
	defaultAddr := envOrDefault("ADDR", ":8080")

	result := config{}
	flag.StringVar(&result.addr, "addr", defaultAddr, "server address, default: "+defaultAddr)
	flag.Parse()
	return result
}

type App struct {
	config config
}

func NewApp() *App {
	return &App{
		config: newConfigFromFlags(),
	}
}

func (a App) Run() error {
	stocksHandler := handlers.NewStocksHandler(services.NewStocksService(repository.NewDumbRepo()))
	http.HandleFunc("/stock/info", stocksHandler.Handle)

	return http.ListenAndServe(a.config.addr, nil)
}
