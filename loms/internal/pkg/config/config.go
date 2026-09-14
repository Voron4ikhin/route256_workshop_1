package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Addr        string
	ProductMock bool
}

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
		defaultAddr        = envOrDefault("LOMS_ADD_LOC", ":8080")
		defaultProductMock = envOrDefaultBool("PRODUCT_MOCK", true) // TODO: вернуть false, когда product-сервис снова станет доступен
	)

	result := Config{}
	flag.StringVar(&result.Addr, "addr", defaultAddr, "server address, default: "+defaultAddr)
	flag.BoolVar(&result.ProductMock, "product_mock", defaultProductMock, "use in-memory mock instead of the real product client")
	flag.Parse()
	return result
}
