package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr     string
	GRPCAddr string
}

func envOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func NewFromFlags() Config {
	defaultAddr := envOrDefault("LOMS_ADD_LOC", ":8080")
	defaultGRPCAddr := envOrDefault("LOMS_GRPC_ADDR_LOC", ":50051")

	result := Config{}
	flag.StringVar(&result.Addr, "addr", defaultAddr, "server address, default: "+defaultAddr)
	flag.StringVar(&result.GRPCAddr, "grpc-addr", defaultGRPCAddr, "gRPC server address, default: "+defaultGRPCAddr)
	flag.Parse()
	return result
}
