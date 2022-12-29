package to_broker

import (
	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/hexy-dev/spacebox-crawler/types"
)

func (tb ToBroker) MapBlock(block *types.Block, totalGas uint64) model.Block {
	return model.Block{
		Height:          block.Height,
		Hash:            block.Hash,
		TxNum:           block.TxNum,
		TotalGas:        totalGas,
		ProposerAddress: block.ProposerAddress,
		Timestamp:       block.Timestamp,
	}
}
