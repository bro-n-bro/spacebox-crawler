package staking

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/hexy-dev/spacebox-crawler/modules/staking/utils"
	tb "github.com/hexy-dev/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/hexy-dev/spacebox-crawler/types"
)

func (m *Module) HandleGenesis(ctx context.Context, doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// Read the genesis state
	var genState stakingtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[stakingtypes.ModuleName], &genState)
	if err != nil {
		return err
	}

	// Save the params
	err = m.publishParams(ctx, doc.InitialHeight, genState.Params)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis params: %s", err)
	}

	// Parse genesis transactions
	err = parseGenesisTransactions(ctx, doc, appState, m.cdc, m.tbM, m.broker)
	if err != nil {
		return fmt.Errorf("error while storing genesis transactions: %s", err)
	}

	// Save the validators
	err = m.publishValidators(ctx, doc, genState.Validators)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validators: %s", err)
	}

	// Save the delegations
	err = m.publishDelegations(ctx, doc, genState)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis delegations: %s", err)
	}

	// Save the unbonding delegations
	err = m.publishUnbondingDelegations(ctx, doc, genState)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis unbonding delegations: %s", err)
	}

	// Save the re-delegations
	err = m.publishRedelegations(ctx, doc, genState)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis redelegations: %s", err)
	}

	// FIXME: dead code?
	// Save the description
	// err = saveValidatorDescription(doc, genState.Validators)
	// if err != nil {
	//	return fmt.Errorf("error while storing staking genesis validator descriptions: %s", err)
	// }

	// FIXME: does it needed?
	// err = publishValidatorsCommissions(doc.InitialHeight, genState.Validators)
	// if err != nil {
	//	return fmt.Errorf("error while storing staking genesis validators commissions: %s", err)
	// }

	return nil
}

func parseGenesisTransactions(ctx context.Context, doc *tmtypes.GenesisDoc,
	appState map[string]json.RawMessage, cdc codec.Codec, mapper tb.ToBroker, broker broker) error {

	var genUtilState genutiltypes.GenesisState
	err := cdc.UnmarshalJSON(appState[genutiltypes.ModuleName], &genUtilState)
	if err != nil {
		return err
	}

	for _, genTxBz := range genUtilState.GetGenTxs() {
		// Unmarshal the transaction
		var genTx tx.Tx
		if err = cdc.UnmarshalJSON(genTxBz, &genTx); err != nil {
			return err
		}

		for _, msg := range genTx.GetMsgs() {
			// Handle the message properly
			createValMsg, ok := msg.(*stakingtypes.MsgCreateValidator)
			if !ok {
				continue
			}

			err = utils.StoreValidatorFromMsgCreateValidator(ctx, doc.InitialHeight, createValMsg, cdc, mapper, broker)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// publishParams saves the given params to the broker.
func (m *Module) publishParams(ctx context.Context, height int64, params stakingtypes.Params) error {
	var commissionRate float64
	if !params.MinCommissionRate.IsNil() {
		commissionRate = params.MinCommissionRate.MustFloat64()
	}

	modelParams := model.NewStakingParams(height, params.MaxValidators, params.MaxEntries, params.HistoricalEntries,
		params.BondDenom, commissionRate, params.UnbondingTime)

	// TODO: test it
	err := m.broker.PublishStakingParams(ctx, modelParams)
	if err != nil {
		return err
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// publishValidators publishes the validators data present inside the given genesis state to the broker.
func (m *Module) publishValidators(ctx context.Context, doc *tmtypes.GenesisDoc, validators stakingtypes.Validators) error {
	vals := make([]types.StakingValidator, len(validators))
	for i, val := range validators {
		validator, err := utils.ConvertValidator(m.cdc, val, doc.InitialHeight)
		if err != nil {
			return err
		}

		vals[i] = validator
	}

	// TODO: save to mongo?
	// TODO test it
	if err := utils.PublishValidatorsData(ctx, vals, m.broker); err != nil {
		return err
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// publishDelegations publishes the delegations and account data present inside the given genesis state to the broker.
func (m *Module) publishDelegations(ctx context.Context, doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState) error {
	for _, validator := range genState.Validators {
		tokens := validator.Tokens
		delegatorShares := validator.DelegatorShares

		typesDelegations := findDelegations(genState.Delegations, validator.OperatorAddress)
		for _, delegation := range typesDelegations {

			delegationAmount := sdk.NewDecFromBigInt(tokens.BigInt()).Mul(delegation.Shares).Quo(delegatorShares).TruncateInt()
			// TODO: test it
			acc := model.NewAccount(delegation.DelegatorAddress, doc.InitialHeight)
			if err := m.broker.PublishAccounts(ctx, []model.Account{acc}); err != nil {
				return err
			}

			// TODO: save to mongo?
			// TODO: test it
			modelDel := model.NewDelegation(validator.OperatorAddress, delegation.DelegatorAddress, doc.InitialHeight,
				m.tbM.MapCoin(types.NewCoinFromCdk(sdk.NewCoin(genState.Params.BondDenom, delegationAmount))))

			if err := m.broker.PublishDelegation(ctx, modelDel); err != nil {
				return err
			}
		}
	}

	return nil
}

// findDelegations returns the list of all the delegations that are
// related to the validator having the given validator address
func findDelegations(genData stakingtypes.Delegations, valAddr string) stakingtypes.Delegations {
	var delegations stakingtypes.Delegations
	for _, delegation := range genData {
		if delegation.ValidatorAddress == valAddr {
			delegations = append(delegations, delegation)
		}
	}
	return delegations
}

// --------------------------------------------------------------------------------------------------------------------

// publishUnbondingDelegations publishes the unbonding delegations data present inside the given genesis state to the broker.
func (m *Module) publishUnbondingDelegations(ctx context.Context, doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState) error {
	for _, validator := range genState.Validators {
		valUD := findUnbondingDelegations(genState.UnbondingDelegations, validator.OperatorAddress)
		for _, ud := range valUD {
			for _, entry := range ud.Entries {
				del := model.NewUnbondingDelegation(
					doc.InitialHeight,
					ud.DelegatorAddress,
					validator.OperatorAddress,
					m.tbM.MapCoin(types.NewCoinFromCdk(sdk.NewCoin(genState.Params.BondDenom, entry.InitialBalance))),
					entry.CompletionTime,
				)

				// TODO: test it
				if err := m.broker.PublishUnbondingDelegation(ctx, del); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// findUnbondingDelegations returns the list of all the unbonding delegations
// that are related to the validator having the given validator address
func findUnbondingDelegations(genData stakingtypes.UnbondingDelegations, valAddr string) stakingtypes.UnbondingDelegations {
	unbondingDelegations := make(stakingtypes.UnbondingDelegations, 0)
	for _, unbondingDelegation := range genData {
		if unbondingDelegation.ValidatorAddress == valAddr {
			unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
		}
	}
	return unbondingDelegations
}

// --------------------------------------------------------------------------------------------------------------------

// publishRedelegations publishes the redelegations data present inside the given genesis state to the broker.
func (m *Module) publishRedelegations(ctx context.Context, doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState) error {
	for _, genRedelegation := range genState.Redelegations {
		for _, entry := range genRedelegation.Entries {
			// TODO: save to mongo?
			// TODO: test it
			redelegation := model.NewRedelegation(
				doc.InitialHeight,
				genRedelegation.DelegatorAddress,
				genRedelegation.ValidatorSrcAddress,
				genRedelegation.ValidatorDstAddress,
				m.tbM.MapCoin(types.NewCoinFromCdk(sdk.NewCoin(genState.Params.BondDenom, entry.InitialBalance))),
				entry.CompletionTime,
			)

			if err := m.broker.PublishRedelegation(ctx, redelegation); err != nil {
				return err
			}
		}
	}

	return nil
}
