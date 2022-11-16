package types

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Module interface {
	Name() string
}

type BlockModule interface {
	// HandleBlock allows to handle a single block.
	// For convenience of use, all the transactions present inside the given block
	// and the currently used database will be passed as well.
	// For each transaction present inside the block, HandleTx will be called as well.
	// NOTE. The returned error will be logged using the logging.LogBlockError method. All other modules' handlers
	// will still be called.
	HandleBlock(context.Context, *Block, *tmctypes.ResultValidators) error
}

type TransactionModule interface {
	// HandleTx handles a single transaction.
	// For each message present inside the transaction, HandleMsg will be called as well.
	// NOTE. The returned error will be logged using the logging.LogTxError method. All other modules' handlers
	// will still be called.
	HandleTx(ctx context.Context, tx *Tx) error
}

type MessageModule interface {
	// HandleMessage handles a single message.
	// For convenience of usa, the index of the message inside the transaction and the transaction itself
	// are passed as well.
	// NOTE. The returned error will be logged using the logging.LogMsgError method. All other modules' handlers
	// will still be called.
	HandleMessage(ctx context.Context, index int, msg sdk.Msg, tx *Tx) error
}

type GenesisModule interface {
	// HandleGenesis allows to handle the genesis state.
	// For convenience of use, the already-unmarshalled AppState is provided along with the full GenesisDoc.
	// NOTE. The returned error will be logged using the logging.LogGenesisError method. All other modules' handlers
	// will still be called.
	HandleGenesis(ctx context.Context, doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error
}
