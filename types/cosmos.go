package types

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type (
	PubKey interface {
		Bytes() []byte
	}

	Block struct {
		Timestamp       time.Time
		Hash            string
		ProposerAddress string
		Evidence        tmtypes.EvidenceData
		TxNum           int
		TotalGas        uint64
		Height          int64
	}

	// Txs - slice of transactions
	Txs []*Tx

	// Tx represents an already existing blockchain transaction
	Tx struct {
		*sdktx.Tx
		*sdk.TxResponse
		Signer string
	}

	Validators []*Validator

	Validator struct {
		ConsAddr   string
		ConsPubkey string
	}

	MessageStruct struct {
	}
)

func NewBlock(height int64, hash, proposerAddress string, txNum int, totalGas uint64, timestamp time.Time,
	evidence tmtypes.EvidenceData) *Block {
	return &Block{
		Height:          height,
		Hash:            hash,
		TxNum:           txNum,
		ProposerAddress: proposerAddress,
		Timestamp:       timestamp,
		Evidence:        evidence,
		TotalGas:        totalGas,
	}
}

// NewBlockFromTmBlock builds a new Block instance from a given ResultBlock object
func NewBlockFromTmBlock(blk *tmctypes.ResultBlock, totalGas uint64) *Block {
	return NewBlock(
		blk.Block.Height,
		blk.Block.Hash().String(),
		sdk.ConsAddress(blk.Block.ProposerAddress).String(),
		len(blk.Block.Txs),
		totalGas,
		blk.Block.Time,
		blk.Block.Evidence,
	)
}

func NewTxsFromTmTxs(txs []*sdktx.GetTxResponse, cdc codec.Codec) Txs {
	res := make(Txs, len(txs))
	for i, tx := range txs {
		var signer string
		if tx.Tx.AuthInfo != nil {
			if len(tx.Tx.AuthInfo.SignerInfos) > 0 {
				if tx.TxResponse.TxHash == "46798FFA86453A448D0FB0484F5345317F6DA6B2715A769EF5981FE5897A8648" {
					println("@")
				}
				var pk cryptotypes.PubKey
				if err := cdc.UnpackAny(tx.Tx.AuthInfo.SignerInfos[0].PublicKey, &pk); err == nil {
					signer, _ = ConvertPubKeyToBech32String(pk)
				}
				// hash46798FFA86453A448D0FB0484F5345317F6DA6B2715A769EF5981FE5897A8648
			}
		}

		res[i] = &Tx{
			Tx:         tx.Tx,
			TxResponse: tx.TxResponse,
			Signer:     signer,
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
		ConsPubkey: consPubKey,
	}
}

func ConvertPubKeyToBech32String(pubKey cryptotypes.PubKey) (string, error) {
	return bech32.ConvertAndEncode("cosmos", pubKey.Bytes()) // TODO
}

// ConvertValidatorPubKeyToBech32String converts the given pubKey to Bech32 string
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

// Successful tells whether this tx is successful or not
func (tx Tx) Successful() bool {
	return tx.TxResponse.Code == 0
}
