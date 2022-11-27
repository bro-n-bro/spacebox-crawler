package rpc

type Config struct {
	Host      string `env:"RPC_URL" envDefault:"http://localhost:26657"`
	WSEnabled bool   `env:"WS_ENABLED"`
}
