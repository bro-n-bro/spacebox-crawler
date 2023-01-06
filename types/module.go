package types

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Module interface {
	// Name base implementation of Module interface.
	Name() string
}

type BlockHandler interface {
	Module
	// HandleBlock handles a single block in blockchain.
	HandleBlock(ctx context.Context, block *Block) error
}

type TransactionHandler interface {
	Module
	// HandleTx handles a single transaction of block.
	HandleTx(ctx context.Context, tx *Tx) error
}

type MessageHandler interface {
	Module
	// HandleMessage handles a single message of transaction.
	HandleMessage(ctx context.Context, index int, msg sdk.Msg, tx *Tx) error
}

type ValidatorsHandler interface {
	Module
	// ValidatorsHandler of all validators in blockchain.
	ValidatorsHandler(ctx context.Context, vals *tmctypes.ResultValidators) error
}

type GenesisHandler interface {
	Module
	// HandleGenesis handles a genesis state.
	HandleGenesis(ctx context.Context, doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error
}
