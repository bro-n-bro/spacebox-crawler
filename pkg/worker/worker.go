package worker

//
//import (
//	"context"
//
//	"bro-n-bro-osmosis/internal/rep"
//	"bro-n-bro-osmosis/types"
//)
//
//type Worker struct {
//	Broker     rep.Broker
//	RPCClient  rep.RPCClient
//	GrpcClient rep.GrpcClient
//	Modules    []types.Module
//
//	cancel func()
//}
//
//func New(b rep.Broker, rpcCli rep.RPCClient, grpcCli rep.GrpcClient, modules []types.Module) *Worker {
//	return &Worker{
//		Broker:     b,
//		RPCClient:  rpcCli,
//		GrpcClient: grpcCli,
//		Modules:    modules,
//	}
//
//}
//
//func (w *Worker) Start(ctx context.Context) error {
//	ctx, cancel := context.WithCancel(ctx)
//
//}
//
//func (w *Worker) Stop(ctx context.Context) error {
//
//}
