package worker

import "bro-n-bro-osmosis/types"

var (
	transactionHandlers []types.TransactionModule
	blockHandlers       []types.BlockModule
	genesisHandlers     []types.GenesisModule
	messageHandlers     []types.MessageModule
)

func (w *Worker) fillModules() {
	txModules := make([]types.TransactionModule, 0)
	msgModules := make([]types.MessageModule, 0)
	genModules := make([]types.GenesisModule, 0)
	blockModules := make([]types.BlockModule, 0)

	for _, module := range w.modules {
		if tI, ok := module.(types.TransactionModule); ok {
			txModules = append(txModules, tI)
		}
		if bI, ok := module.(types.BlockModule); ok {
			blockModules = append(blockModules, bI)
		}
		if gI, ok := module.(types.GenesisModule); ok {
			genModules = append(genModules, gI)
		}
		if mI, ok := module.(types.MessageModule); ok {
			msgModules = append(msgModules, mI)
		}
	}

	transactionHandlers = make([]types.TransactionModule, len(txModules))
	copy(transactionHandlers, txModules)

	blockHandlers = make([]types.BlockModule, len(blockModules))
	copy(blockHandlers, blockModules)

	genesisHandlers = make([]types.GenesisModule, len(genModules))
	copy(genesisHandlers, genModules)

	messageHandlers = make([]types.MessageModule, len(msgModules))
	copy(messageHandlers, msgModules)

}
