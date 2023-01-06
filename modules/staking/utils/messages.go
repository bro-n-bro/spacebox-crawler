package utils

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	grpcClient "github.com/hexy-dev/spacebox-crawler/client/grpc"
	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/hexy-dev/spacebox-crawler/types"
	"github.com/hexy-dev/spacebox/broker/model"
)

// StoreValidatorFromMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func StoreValidatorFromMsgCreateValidator(
	ctx context.Context,
	height int64,
	msg *stakingtypes.MsgCreateValidator,
	cdc codec.Codec,
	mapper tb.ToBroker,
	broker interface {
		PublishAccounts(ctx context.Context, accounts []model.Account) error // FIXME: auth module
		PublishValidators(ctx context.Context, vals []model.Validator) error
		PublishValidatorsInfo(ctx context.Context, infos []model.ValidatorInfo) error
		PublishUnbondingDelegation(ctx context.Context, ud model.UnbondingDelegation) error
		PublishDelegation(ctx context.Context, d model.Delegation) error
	},
) error {

	var pubKey cryptotypes.PubKey
	err := cdc.UnpackAny(msg.Pubkey, &pubKey)
	if err != nil {
		return err
	}

	operatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return err
	}

	stakingValidator, err := stakingtypes.NewValidator(operatorAddr, pubKey, msg.Description)
	if err != nil {
		return err
	}

	validator, err := ConvertValidator(cdc, stakingValidator, height)
	if err != nil {
		return err
	}

	// TODO: does it needed?
	// desc, err := ConvertValidatorDescription(msg.ValidatorAddress, msg.Description, height)
	// if err != nil {
	//	return err
	// }

	// TODO: save to mongo?
	// TODO: test it
	if err = PublishValidatorsData(ctx, []types.StakingValidator{validator}, broker); err != nil {
		return err
	}

	// TODO: save to mongo?
	// TODO: test it
	// Save the first self-delegation
	if err = broker.PublishDelegation(ctx, model.NewDelegation(
		msg.ValidatorAddress,
		msg.DelegatorAddress,
		height,
		mapper.MapCoin(types.NewCoinFromCdk(msg.Value)),
	)); err != nil {
		return err
	}

	// FIXME: does it needed?
	// Save the description
	// err = broker.PublishValidatorDescription(desc)
	// if err != nil {
	//	return err
	// }
	//

	// Save the commission
	// err = broker.publishValidatorCommission(types.NewValidatorCommission(
	//	msg.ValidatorAddress,
	//	&msg.Commission.Rate,
	//	&msg.MinSelfDelegation,
	//	height,
	// ))
	return err
}

// StoreDelegationFromMessage handles a MsgDelegate and saves the delegation inside the database
func StoreDelegationFromMessage(ctx context.Context, tx *types.Tx, msg *stakingtypes.MsgDelegate,
	stakingClient stakingtypes.QueryClient, mapper tb.ToBroker, broker interface {
		PublishDelegation(ctx context.Context, d model.Delegation) error
		PublishDelegationMessage(ctx context.Context, dm model.DelegationMessage) error
	}) error {

	header := grpcClient.GetHeightRequestHeader(tx.Height)
	res, err := stakingClient.Delegation(
		context.Background(),
		&stakingtypes.QueryDelegationRequest{
			DelegatorAddr: msg.DelegatorAddress,
			ValidatorAddr: msg.ValidatorAddress,
		},
		header,
	)
	if err != nil {
		return err
	}

	// TODO: test it
	d := model.NewDelegation(
		res.DelegationResponse.Delegation.ValidatorAddress,
		res.DelegationResponse.Delegation.DelegatorAddress,
		tx.Height,
		mapper.MapCoin(types.NewCoinFromCdk(res.DelegationResponse.Balance)),
	)

	if err = broker.PublishDelegation(ctx, d); err != nil {
		return err
	}

	dm := model.NewDelegationMessage(
		res.DelegationResponse.Delegation.DelegatorAddress,
		res.DelegationResponse.Delegation.ValidatorAddress,
		tx.TxHash,
		tx.Height,
		mapper.MapCoin(types.NewCoinFromCdk(res.DelegationResponse.Balance)),
	)

	if err = broker.PublishDelegationMessage(ctx, dm); err != nil {
		return err
	}

	return nil
}
