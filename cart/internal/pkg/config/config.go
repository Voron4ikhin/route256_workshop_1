package config

import (
	"flag"
	"os"
	"strconv"
)

// Config is shared by every transport (http, grpc, ...) the app exposes.
type Config struct {
	Addr        string
	LomsAddr    string
	ProductAddr string
	ProductMock bool
}

// envOrDefault returns the value of the env var if set, otherwise def.
// Flags below use this as their default, so the precedence is: flag > env > hardcoded default.
func envOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func envOrDefaultBool(key string, def bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func NewFromFlags() Config {
	var (
		defaultAddr        = envOrDefault("ADDR", ":8080")
		defaultLomsAddr    = envOrDefault("LOMS_ADDR", "http://loms:8080")
		defaultProductAddr = envOrDefault("PRODUCT_ADDR", "http://route256.pavl.uk:8080") // это не работает
		defaultProductMock = envOrDefaultBool("PRODUCT_MOCK", true)                       // TODO: вернуть false, когда product-сервис снова станет доступен
	)

	result := Config{}
	flag.StringVar(&result.Addr, "addr", defaultAddr, "server address, default: "+defaultAddr)
	flag.StringVar(&result.LomsAddr, "loms_addr", defaultLomsAddr, "loms server address, default: "+defaultLomsAddr)
	flag.StringVar(&result.ProductAddr, "product_addr", defaultProductAddr, "product server address, default: "+defaultProductAddr)
	flag.BoolVar(&result.ProductMock, "product_mock", defaultProductMock, "use in-memory mock instead of the real product client")
	flag.Parse()
	return result
}
