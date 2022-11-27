package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type (
	Block struct {
		Height          int64
		Hash            string
		TxNum           int
		ProposerAddress string
		Timestamp       time.Time
		Evidence        tmtypes.EvidenceData
	}

	// Txs - slice of transactions
	Txs []*Tx

	// Tx represents an already existing blockchain transaction
	Tx struct {
		*sdktx.Tx
		*sdk.TxResponse
	}

	Validators []*Validator

	Validator struct {
		ConsAddr   string
		ConsPubKey string
	}
)

func NewBlock(height int64, hash, proposerAddress string, txNum int, timestamp time.Time, evidence tmtypes.EvidenceData) *Block {
	return &Block{
		Height:          height,
		Hash:            hash,
		TxNum:           txNum,
		ProposerAddress: proposerAddress,
		Timestamp:       timestamp,
		Evidence:        evidence,
	}
}

// NewBlockFromTmBlock builds a new Block instance from a given ResultBlock object
func NewBlockFromTmBlock(blk *tmctypes.ResultBlock) *Block {
	return NewBlock(
		blk.Block.Height,
		blk.Block.Hash().String(),
		sdk.ConsAddress(blk.Block.ProposerAddress).String(),
		len(blk.Block.Txs),
		blk.Block.Time,
		blk.Block.Evidence,
	)
}
func NewTxsFromTmTxs(txs []*sdktx.GetTxResponse) Txs {
	res := make(Txs, len(txs))
	for i, tx := range txs {
		res[i] = &Tx{
			Tx:         tx.Tx,
			TxResponse: tx.TxResponse,
		}
	}
	return res
}

func NewValidatorsFromTmValidator(tmVals *tmctypes.ResultValidators) Validators {
	res := make(Validators, 0, len(tmVals.Validators))
	for _, val := range tmVals.Validators {
		consAddr := sdk.ConsAddress(val.Address).String()
		consPubKey, err := ConvertValidatorPubKeyToBech32String(val.PubKey)
		if err == nil {
			res = append(res, NewValidator(consAddr, consPubKey))
		}
	}
	return res
}

// NewValidator allows to build a new Validator instance
func NewValidator(consAddr string, consPubKey string) *Validator {
	return &Validator{
		ConsAddr:   consAddr,
		ConsPubKey: consPubKey,
	}
}

// ConvertValidatorPubKeyToBech32String converts the given pubKey to a Bech32 string
func ConvertValidatorPubKeyToBech32String(pubKey tmcrypto.PubKey) (string, error) {
	bech32Prefix := sdk.GetConfig().GetBech32ConsensusPubPrefix()
	return bech32.ConvertAndEncode(bech32Prefix, pubKey.Bytes())
}

// TotalGas calculates and returns total used gas of all transactions
func (txs Txs) TotalGas() (totalGas uint64) {
	for _, tx := range txs {
		totalGas += uint64(tx.GasUsed)
	}
	return totalGas
}

// FindEventByType searches inside the given tx events for the message having the specified index, in order
// to find the event having the given type, and returns it.
// If no such event is found, returns an error instead.
func (tx Tx) FindEventByType(index int, eventType string) (sdk.StringEvent, error) {
	for _, ev := range tx.Logs[index].Events {
		if ev.Type == eventType {
			return ev, nil
		}
	}

	return sdk.StringEvent{}, fmt.Errorf("no %s event found inside tx with hash %s", eventType, tx.TxHash)
}

// FindAttributeByKey searches inside the specified event of the given tx to find the attribute having the given key.
// If the specified event does not contain a such attribute, returns an error instead.
func (tx Tx) FindAttributeByKey(event sdk.StringEvent, attrKey string) (string, error) {
	for _, attr := range event.Attributes {
		if attr.Key == attrKey {
			return attr.Value, nil
		}
	}

	return "", fmt.Errorf("no event with attribute %s found inside tx with hash %s", attrKey, tx.TxHash)
}
