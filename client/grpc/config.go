package grpc

type Config struct {
	Host string `env:"GRPC_URL" envDefault:"http://localhost:9090"`
}
