// nolint
package liquidity

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// liquidity module sentinel errors
var (
	ErrNumOfReserveCoin         = sdkerrors.Register(ModuleName, 5, "invalid number of reserve coin")
	ErrInvalidPoolCreatorAddr   = sdkerrors.Register(ModuleName, 15, "invalid pool creator address")
	ErrInvalidDepositorAddr     = sdkerrors.Register(ModuleName, 16, "invalid pool depositor address")
	ErrInvalidWithdrawerAddr    = sdkerrors.Register(ModuleName, 17, "invalid pool withdrawer address")
	ErrInvalidSwapRequesterAddr = sdkerrors.Register(ModuleName, 18, "invalid pool swap requester address")
	ErrBadPoolCoinAmount        = sdkerrors.Register(ModuleName, 19, "invalid pool coin amount")
	ErrBadDepositCoinsAmount    = sdkerrors.Register(ModuleName, 20, "invalid deposit coins amount")
	ErrBadOfferCoinAmount       = sdkerrors.Register(ModuleName, 21, "invalid offer coin amount")
	ErrBadOrderPrice            = sdkerrors.Register(ModuleName, 23, "invalid order price")
	ErrLessThanMinOfferAmount   = sdkerrors.Register(ModuleName, 34, "offer amount should be over 100 micro")
	ErrBadPoolTypeID            = sdkerrors.Register(ModuleName, 37, "invalid index of the pool type")
)
