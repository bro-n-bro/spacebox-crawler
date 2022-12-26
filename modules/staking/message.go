package staking

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"bro-n-bro-osmosis/internal/rep"
	stakingutils "bro-n-bro-osmosis/modules/staking/utils"
	tb "bro-n-bro-osmosis/pkg/mapper/to_broker"
	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleMessage(ctx context.Context, index int, cosmosMsg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch msg := cosmosMsg.(type) {
	case *stakingtypes.MsgCreateValidator:
		return handleMsgCreateValidator(ctx, tx.Height, msg, m.cdc, m.broker, m.tbM)

	case *stakingtypes.MsgEditValidator:
		return handleEditValidator(tx.Height, msg)

	case *stakingtypes.MsgDelegate:
		return stakingutils.StoreDelegationFromMessage(ctx, tx, msg, m.client.StakingQueryClient, m.broker, m.tbM)

	case *stakingtypes.MsgBeginRedelegate:
		return handleMsgBeginRedelegate(ctx, tx, index, msg, m.client.StakingQueryClient, m.broker, m.tbM)

	case *stakingtypes.MsgUndelegate:
		return handleMsgUndelegate(ctx, tx, index, msg, m.client.StakingQueryClient, m.broker, m.tbM)
	}

	return nil
}

// handleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func handleMsgCreateValidator(ctx context.Context, height int64, msg *stakingtypes.MsgCreateValidator, cdc codec.Codec,
	broker rep.Broker, mapper tb.ToBroker) error {

	err := stakingutils.StoreValidatorFromMsgCreateValidator(ctx, height, msg, cdc, broker, mapper)
	if err != nil {
		return err
	}

	// Save validator description
	_, err = stakingutils.ConvertValidatorDescription(
		msg.ValidatorAddress,
		msg.Description,
		height,
	)
	if err != nil {
		return err
	}

	//err = db.SaveValidatorDescription(description)
	//if err != nil {
	//	return err
	//}

	// Save validator commission
	//return db.SaveValidatorCommission(types.NewValidatorCommission(
	//	msg.ValidatorAddress,
	//	&msg.Commission.Rate,
	//	&msg.MinSelfDelegation,
	//	height,
	//))
	return nil
}

// handleEditValidator handles MsgEditValidator utils, updating the validator info and commission
func handleEditValidator(height int64, msg *stakingtypes.MsgEditValidator) error {
	// Save validator commission
	//err := db.SaveValidatorCommission(types.NewValidatorCommission(
	//	msg.ValidatorAddress,
	//	msg.CommissionRate,
	//	msg.MinSelfDelegation,
	//	height,
	//))
	//if err != nil {
	//	return err
	//}

	// Save validator description
	_, err := stakingutils.ConvertValidatorDescription(
		msg.ValidatorAddress,
		msg.Description,
		height,
	)
	if err != nil {
		return err
	}

	return nil
}

// handleMsgBeginRedelegate handles a MsgBeginRedelegate storing the data inside the database
func handleMsgBeginRedelegate(ctx context.Context, tx *types.Tx, index int, msg *stakingtypes.MsgBeginRedelegate,
	client stakingtypes.QueryClient, broker rep.Broker, mapper tb.ToBroker,
) error {
	_, err := stakingutils.StoreRedelegationFromMessage(ctx, tx, index, msg, broker, mapper)
	if err != nil {
		return err
	}

	// Update the current delegations
	return stakingutils.UpdateDelegationsAndReplaceExisting(ctx, tx.Height, msg.DelegatorAddress, client, broker, mapper)
}

// handleMsgUndelegate handles a MsgUndelegate storing the data inside the database
func handleMsgUndelegate(ctx context.Context, tx *types.Tx, index int, msg *stakingtypes.MsgUndelegate,
	stakingClient stakingtypes.QueryClient, broker rep.Broker, mapper tb.ToBroker,
) error {
	_, err := stakingutils.StoreUnbondingDelegationFromMessage(ctx, tx, index, msg, broker, mapper)
	if err != nil {
		return err
	}

	// Update the current delegations
	return stakingutils.UpdateDelegationsAndReplaceExisting(ctx, tx.Height, msg.DelegatorAddress, stakingClient, broker, mapper)
}
