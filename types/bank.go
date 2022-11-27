package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type (
	// AccountBalance represents the balance of an account at a given height
	AccountBalance struct {
		Address string
		Balance sdk.Coins
		Height  int64
	}

	TotalSupply struct {
		Height int64
		Coins  Coins
	}

	MsgSend struct {
		Coins       Coins
		AddressFrom string
		AddressTo   string
		TxHash      string
		Height      int64
	}
)

// NewAccountBalance allows to build a new AccountBalance instance
func NewAccountBalance(address string, balance sdk.Coins, height int64) AccountBalance {
	return AccountBalance{
		Address: address,
		Balance: balance,
		Height:  height,
	}
}

// NewTotalSupply returns the new TotalSupply instance
func NewTotalSupply(height int64, coins Coins) TotalSupply {
	return TotalSupply{
		Height: height,
		Coins:  coins,
	}
}

func NewMsgSend(coins Coins, height int64, addressFrom, addressTo, txHash string) MsgSend {
	return MsgSend{
		Coins:       coins,
		AddressFrom: addressFrom,
		AddressTo:   addressTo,
		TxHash:      txHash,
		Height:      height,
	}
}
