package types

import (
	"time"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// ValidatorSigningInfo contains the signing info of a validator at a given height
type ValidatorSigningInfo struct {
	JailedUntil         time.Time
	ValidatorAddress    string
	StartHeight         int64
	IndexOffset         int64
	MissedBlocksCounter int64
	Height              int64
	Tombstoned          bool
}

// Equal tells whether v and w represent the same rows
func (v ValidatorSigningInfo) Equal(w ValidatorSigningInfo) bool {
	return v.ValidatorAddress == w.ValidatorAddress &&
		v.StartHeight == w.StartHeight &&
		v.IndexOffset == w.IndexOffset &&
		v.JailedUntil.Equal(w.JailedUntil) &&
		v.Tombstoned == w.Tombstoned &&
		v.MissedBlocksCounter == w.MissedBlocksCounter &&
		v.Height == w.Height
}

// ValidatorSigningInfo allows to build a new ValidatorSigningInfo
func NewValidatorSigningInfo(validatorAddress string, height, startHeight, indexOffset, missedBlocksCounter int64,
	jailedUntil time.Time, tombstoned bool) ValidatorSigningInfo {
	return ValidatorSigningInfo{
		ValidatorAddress:    validatorAddress,
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		JailedUntil:         jailedUntil,
		Tombstoned:          tombstoned,
		MissedBlocksCounter: missedBlocksCounter,
		Height:              height,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// SlashingParams represents the parameters of the slashing module at a given height
type SlashingParams struct {
	slashingtypes.Params
	Height int64
}

// NewSlashingParams allows to build a new SlashingParams instance
func NewSlashingParams(params slashingtypes.Params, height int64) SlashingParams {
	return SlashingParams{
		Params: params,
		Height: height,
	}
}
