package worker

import "github.com/bro-n-bro/spacebox-crawler/types"

var (
	transactionHandlers  []types.TransactionHandler
	blockHandlers        []types.BlockHandler
	genesisHandlers      []types.GenesisHandler
	messageHandlers      []types.MessageHandler
	validatorsHandlers   []types.ValidatorsHandler
	beginBlockerHandlers []types.BeginBlockerHandler
	endBlockerHandlers   []types.EndBlockerHandler
)

// fillModules fills the module handlers.
func (w *Worker) fillModules() {
	for _, module := range w.modules {
		if tI, ok := module.(types.TransactionHandler); ok {
			transactionHandlers = append(transactionHandlers, tI)
		}
		if bI, ok := module.(types.BlockHandler); ok {
			blockHandlers = append(blockHandlers, bI)
		}
		if gI, ok := module.(types.GenesisHandler); ok {
			genesisHandlers = append(genesisHandlers, gI)
		}
		if mI, ok := module.(types.MessageHandler); ok {
			messageHandlers = append(messageHandlers, mI)
		}
		if vI, ok := module.(types.ValidatorsHandler); ok {
			validatorsHandlers = append(validatorsHandlers, vI)
		}
		if bbI, ok := module.(types.BeginBlockerHandler); ok {
			beginBlockerHandlers = append(beginBlockerHandlers, bbI)
		}
		if ebI, ok := module.(types.EndBlockerHandler); ok {
			endBlockerHandlers = append(endBlockerHandlers, ebI)
		}
	}
}
