package worker

import "github.com/bro-n-bro/spacebox-crawler/types"

var (
	transactionHandlers []types.TransactionHandler
	blockHandlers       []types.BlockHandler
	genesisHandlers     []types.GenesisHandler
	messageHandlers     []types.MessageHandler
	validatorsHandlers  []types.ValidatorsHandler
)

func (w *Worker) fillModules() {
	txModules := make([]types.TransactionHandler, 0)
	msgModules := make([]types.MessageHandler, 0)
	genModules := make([]types.GenesisHandler, 0)
	blockModules := make([]types.BlockHandler, 0)
	validatorsModules := make([]types.ValidatorsHandler, 0)

	for _, module := range w.modules {
		if tI, ok := module.(types.TransactionHandler); ok {
			txModules = append(txModules, tI)
		}
		if bI, ok := module.(types.BlockHandler); ok {
			blockModules = append(blockModules, bI)
		}
		if gI, ok := module.(types.GenesisHandler); ok {
			genModules = append(genModules, gI)
		}
		if mI, ok := module.(types.MessageHandler); ok {
			msgModules = append(msgModules, mI)
		}
		if vI, ok := module.(types.ValidatorsHandler); ok {
			validatorsModules = append(validatorsModules, vI)
		}
	}

	transactionHandlers = make([]types.TransactionHandler, len(txModules))
	copy(transactionHandlers, txModules)

	blockHandlers = make([]types.BlockHandler, len(blockModules))
	copy(blockHandlers, blockModules)

	genesisHandlers = make([]types.GenesisHandler, len(genModules))
	copy(genesisHandlers, genModules)

	messageHandlers = make([]types.MessageHandler, len(msgModules))
	copy(messageHandlers, msgModules)

	validatorsHandlers = make([]types.ValidatorsHandler, len(validatorsModules))
	copy(validatorsHandlers, validatorsModules)
}
