package staking

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func (m *Module) HandleGenesis(ctx context.Context, doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// Read the genesis state
	var genState stakingtypes.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[stakingtypes.ModuleName], &genState); err != nil {
		return err
	}

	// Save the params
	if err := m.publishParams(ctx, doc.InitialHeight, genState.Params); err != nil {
		return fmt.Errorf("error while storing staking genesis params: %w", err)
	}

	// Parse genesis transactions
	if err := m.parseGenesisTransactions(ctx, doc, appState); err != nil {
		return fmt.Errorf("error while storing genesis transactions: %w", err)
	}

	// Save the validators
	if err := m.publishValidators(ctx, doc, genState.Validators); err != nil {
		return fmt.Errorf("error while storing staking genesis validators: %w", err)
	}

	// Save the delegations
	if err := m.publishDelegations(ctx, doc, genState); err != nil {
		return fmt.Errorf("error while storing staking genesis delegations: %w", err)
	}

	// Save the unbonding delegations
	if err := m.publishUnbondingDelegations(ctx, doc, genState); err != nil {
		return fmt.Errorf("error while storing staking genesis unbonding delegations: %w", err)
	}

	// Save the re-delegations
	if err := m.publishRedelegations(ctx, doc, genState); err != nil {
		return fmt.Errorf("error while storing staking genesis redelegations: %w", err)
	}

	// Save the description
	if err := m.publishValidatorDescription(ctx, doc, genState.Validators); err != nil {
		return fmt.Errorf("error while storing staking genesis validator descriptions: %w", err)
	}

	// FIXME: does it needed?
	// if err := publishValidatorsCommissions(doc.InitialHeight, genState.Validators); err != nil {
	//	return fmt.Errorf("error while storing staking genesis validators commissions: %s", err)
	// }

	return nil
}

func (m *Module) parseGenesisTransactions(ctx context.Context, doc *tmtypes.GenesisDoc,
	appState map[string]json.RawMessage) error {

	var genUtilState genutiltypes.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[genutiltypes.ModuleName], &genUtilState); err != nil {
		return err
	}

	for _, genTxBz := range genUtilState.GetGenTxs() {
		// Unmarshal the transaction
		var genTx tx.Tx
		if err := m.cdc.UnmarshalJSON(genTxBz, &genTx); err != nil {
			return err
		}

		for _, msg := range genTx.GetMsgs() {
			// Handle the message properly
			createValMsg, ok := msg.(*stakingtypes.MsgCreateValidator)
			if !ok {
				continue
			}

			if err := m.handleMsgCreateValidator(ctx, doc.InitialHeight, "", 0, createValMsg); err != nil {
				return err
			}
		}
	}

	return nil
}

// publishParams saves the given params to the broker.
func (m *Module) publishParams(ctx context.Context, height int64, params stakingtypes.Params) error {
	var commissionRate float64
	if !params.MinCommissionRate.IsNil() {
		commissionRate = params.MinCommissionRate.MustFloat64()
	}

	// TODO: test it
	return m.broker.PublishStakingParams(ctx, model.StakingParams{
		Height: height,
		Params: model.SParams{
			UnbondingTime:     params.UnbondingTime,
			MaxValidators:     uint64(params.MaxValidators),
			MaxEntries:        uint64(params.MaxEntries),
			HistoricalEntries: uint64(params.HistoricalEntries),
			BondDenom:         params.BondDenom,
			MinCommissionRate: commissionRate,
		},
	})
}

// publishValidators publishes the validators data present inside the given genesis state to the broker.
func (m *Module) publishValidators(ctx context.Context, doc *tmtypes.GenesisDoc, validators stakingtypes.Validators) error {
	vals := make([]types.StakingValidator, len(validators))

	for i, val := range validators {
		validator, err := convertValidator(m.cdc, val, doc.InitialHeight)
		if err != nil {
			return err
		}

		vals[i] = validator
	}

	// TODO: save to mongo?
	// TODO test it
	return m.PublishValidatorsData(ctx, vals)
}

// publishDelegations publishes the delegations and account data present inside the given genesis state to the broker.
func (m *Module) publishDelegations(ctx context.Context, doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState) error {
	prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()

	for _, validator := range genState.Validators {
		tokens := validator.Tokens
		delShares := validator.DelegatorShares

		typesDelegations := findDelegations(genState.Delegations, validator.OperatorAddress)

		for _, del := range typesDelegations {
			if strings.HasPrefix(del.DelegatorAddress, prefix) {
				// TODO: test it
				if err := m.broker.PublishAccount(ctx, model.Account{
					Address: del.DelegatorAddress,
					Height:  doc.InitialHeight,
				}); err != nil {
					return err
				}
			}

			delegationAmount := sdk.NewDecFromBigInt(tokens.BigInt()).Mul(del.Shares).Quo(delShares).TruncateInt()
			// TODO: save to mongo?
			// TODO: test it
			if err := m.broker.PublishDelegation(ctx, model.Delegation{
				OperatorAddress:  validator.OperatorAddress,
				DelegatorAddress: del.DelegatorAddress,
				Height:           doc.InitialHeight,
				Coin: m.tbM.MapCoin(
					types.NewCoinFromCdk(sdk.NewCoin(genState.Params.BondDenom, delegationAmount))),
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

// publishUnbondingDelegations publishes the unbonding delegations data present inside the given genesis state to the broker.
func (m *Module) publishUnbondingDelegations(ctx context.Context, doc *tmtypes.GenesisDoc,
	genState stakingtypes.GenesisState) error {

	var coin types.Coin

	for _, validator := range genState.Validators {
		valUD := findUnbondingDelegations(genState.UnbondingDelegations, validator.OperatorAddress)
		for _, ud := range valUD {
			for _, entry := range ud.Entries {
				coin = types.NewCoinFromCdk(sdk.NewCoin(genState.Params.BondDenom, entry.InitialBalance))
				// TODO: test it
				if err := m.broker.PublishUnbondingDelegation(ctx, model.UnbondingDelegation{
					Height:              doc.InitialHeight,
					DelegatorAddress:    ud.DelegatorAddress,
					ValidatorAddress:    validator.OperatorAddress,
					Coin:                m.tbM.MapCoin(coin),
					CompletionTimestamp: entry.CompletionTime,
				}); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// publishRedelegations publishes the redelegations data present inside the given genesis state to the broker.
func (m *Module) publishRedelegations(ctx context.Context, doc *tmtypes.GenesisDoc,
	genState stakingtypes.GenesisState) error {

	for _, genRedelegation := range genState.Redelegations {
		for _, entry := range genRedelegation.Entries {
			// TODO: save to mongo?
			// TODO: test it
			if err := m.broker.PublishRedelegation(ctx, model.Redelegation{
				Height:              doc.InitialHeight,
				DelegatorAddress:    genRedelegation.DelegatorAddress,
				SrcValidatorAddress: genRedelegation.ValidatorSrcAddress,
				DstValidatorAddress: genRedelegation.ValidatorDstAddress,
				Coin: m.tbM.MapCoin(
					types.NewCoinFromCdk(sdk.NewCoin(genState.Params.BondDenom, entry.InitialBalance))),
				CompletionTime: entry.CompletionTime,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Module) publishValidatorDescription(ctx context.Context, doc *tmtypes.GenesisDoc,
	validators stakingtypes.Validators) error {

	for _, val := range validators {
		if err := m.broker.PublishValidatorDescription(ctx, model.ValidatorDescription{
			OperatorAddress: val.OperatorAddress,
			Moniker:         val.GetMoniker(),
			Identity:        val.Description.Identity,
			Website:         val.Description.Website,
			SecurityContact: val.Description.SecurityContact,
			Details:         val.Description.Details,
			AvatarURL:       "", // TODO:
			Height:          doc.InitialHeight,
		}); err != nil {
			return err
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
