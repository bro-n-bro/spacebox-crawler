package app

import (
	"bro-n-bro-osmosis/adapter/broker"
	"bro-n-bro-osmosis/client/grpc"
	"bro-n-bro-osmosis/client/rpc"
	"bro-n-bro-osmosis/pkg/worker"
)

// CHAIN_PREFIX=prefix # префикс индексируемой сети
// CLICKHOUSE_DB_FOLDER=$HOME/.space-box/clickhouse # место хранения файлов БД
// CLICKHOUSE_DB_HOST=localhost # хост БД
// CLICKHOUSE_DB_PORT=5432 # порт БД
// CLICKHOUSE_DB_NAME=space-box # назване БД
// CLICKHOUSE_USER_NAME=space-box # аккаунт пользователя БД
// CLICKHOUSE_DB_PASSWORD=space-box # пароль пользователя БД
// CLICKHOUSE_SSL_MODE=disable # хз, че это
// HASURA_PORT=8090 # порт хасуры
// HASURA_ADMIN_SECRET=hasura # пароль хасуры
// SPACE_BOX_WORKERS=1 # количество параллельных процессов
// RPC_URL=http://localhost:26657 # RPC API индексируемой сети
// GRPC_URL=http://localhost:9090 # GRPC API индексируемой цепи
// START_HEIGHT=1 # Блок, с которого проводить индексацию
// START_FROM_SNAPSHOT=False # флан на случай старта неполной ноды, а со снапшота
// WS_ENABLED=False # обработка новых блоков
type Config struct {
	ChainPrefix string   `env:"CHAIN_PREFIX"`
	Modules     []string `env:"MODULES" required:"true"`

	RpcConfig    rpc.Config
	GrpcConfig   grpc.Config
	BrokerConfig broker.Config
	WorkerConfig worker.Config
}
