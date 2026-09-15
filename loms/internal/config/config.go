package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr string
}

func envOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func NewFromFlags() Config {
	defaultAddr := envOrDefault("LOMS_ADD_LOC", ":8080")

	result := Config{}
	flag.StringVar(&result.Addr, "addr", defaultAddr, "server address, default: "+defaultAddr)
	flag.Parse()
	return result
}
