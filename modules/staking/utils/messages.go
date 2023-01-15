package utils

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		PublishAccount(ctx context.Context, account model.Account) error // FIXME: auth module
		PublishValidator(ctx context.Context, val model.Validator) error
		PublishValidatorInfo(ctx context.Context, info model.ValidatorInfo) error
		PublishUnbondingDelegation(ctx context.Context, ud model.UnbondingDelegation) error
		PublishDelegation(ctx context.Context, d model.Delegation) error
	},
) error {

	var pubKey cryptotypes.PubKey
	if err := cdc.UnpackAny(msg.Pubkey, &pubKey); err != nil {
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
	if err = broker.PublishDelegation(ctx, model.Delegation{
		OperatorAddress:  msg.ValidatorAddress,
		DelegatorAddress: msg.DelegatorAddress,
		Height:           height,
		Coin:             mapper.MapCoin(types.NewCoinFromCdk(msg.Value)),
	}); err != nil {
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
func StoreDelegationFromMessage(
	ctx context.Context,
	tx *types.Tx,
	msg *stakingtypes.MsgDelegate,
	index int,
	stakingClient stakingtypes.QueryClient,
	mapper tb.ToBroker,
	broker interface {
		PublishDelegation(ctx context.Context, d model.Delegation) error
		PublishDelegationMessage(ctx context.Context, dm model.DelegationMessage) error
	},
) error {

	// TODO: test it
	if err := broker.PublishDelegationMessage(ctx, model.DelegationMessage{
		Delegation: model.Delegation{
			OperatorAddress:  msg.ValidatorAddress,
			DelegatorAddress: msg.DelegatorAddress,
			Coin:             mapper.MapCoin(types.NewCoinFromCdk(msg.Amount)),
			Height:           tx.Height,
		},
		TxHash:   tx.TxHash,
		MsgIndex: int64(index),
	}); err != nil {
		return err
	}

	header := grpcClient.GetHeightRequestHeader(tx.Height)

	respPb, err := stakingClient.Delegation(
		ctx,
		&stakingtypes.QueryDelegationRequest{
			DelegatorAddr: msg.DelegatorAddress,
			ValidatorAddr: msg.ValidatorAddress,
		},
		header,
	)
	if err != nil {
		s, ok := status.FromError(err)
		if !ok || s.Code() != codes.NotFound {
			return err
		}
	}

	var coin model.Coin
	if err == nil {
		coin = mapper.MapCoin(types.NewCoinFromCdk(respPb.DelegationResponse.Balance))
	}

	// TODO: test it
	if err = broker.PublishDelegation(ctx, model.Delegation{
		OperatorAddress:  msg.ValidatorAddress,
		DelegatorAddress: msg.DelegatorAddress,
		Height:           tx.Height,
		Coin:             coin,
	}); err != nil {
		return err
	}

	return nil
}
