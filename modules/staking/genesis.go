package staking

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"bro-n-bro-osmosis/modules/staking/utils"
	"bro-n-bro-osmosis/types"
)

func (m *Module) HandleGenesis(_ context.Context, doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// Read the genesis state
	var genState stakingtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[stakingtypes.ModuleName], &genState)
	if err != nil {
		return err
	}

	// Save the params
	err = saveParams(doc.InitialHeight, genState.Params)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis params: %s", err)
	}

	// Parse genesis transactions
	err = parseGenesisTransactions(doc, appState, m.cdc)
	if err != nil {
		return fmt.Errorf("error while storing genesis transactions: %s", err)
	}

	// Save the validators
	err = saveValidators(doc, genState.Validators, m.cdc)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validators: %s", err)
	}

	// Save the delegations
	err = saveDelegations(doc, genState)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis delegations: %s", err)
	}

	// Save the unbonding delegations
	err = saveUnbondingDelegations(doc, genState)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis unbonding delegations: %s", err)
	}

	// Save the re-delegations
	err = saveRedelegations(doc, genState)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis redelegations: %s", err)
	}

	// Save the description
	err = saveValidatorDescription(doc, genState.Validators)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validator descriptions: %s", err)
	}

	err = saveValidatorsCommissions(doc.InitialHeight, genState.Validators)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validators commissions: %s", err)
	}

	return nil
}

func parseGenesisTransactions(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc codec.Codec) error {
	var genUtilState genutiltypes.GenesisState
	err := cdc.UnmarshalJSON(appState[genutiltypes.ModuleName], &genUtilState)
	if err != nil {
		return err
	}

	for _, genTxBz := range genUtilState.GetGenTxs() {
		// Unmarshal the transaction
		var genTx tx.Tx
		if err := cdc.UnmarshalJSON(genTxBz, &genTx); err != nil {
			return err
		}

		for _, msg := range genTx.GetMsgs() {
			// Handle the message properly
			createValMsg, ok := msg.(*stakingtypes.MsgCreateValidator)
			if !ok {
				continue
			}

			err = utils.StoreValidatorFromMsgCreateValidator(doc.InitialHeight, createValMsg, cdc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// saveParams saves the given params into the database
func saveParams(height int64, params stakingtypes.Params) error {
	//return db.SaveStakingParams(types.NewStakingParams(params, height))
	_ = types.NewStakingParams(params, height)
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// saveValidators stores the validators data present inside the given genesis state
func saveValidators(doc *tmtypes.GenesisDoc, validators stakingtypes.Validators, cdc codec.Codec) error {
	vals := make([]types.StakingValidator, len(validators))
	for i, val := range validators {
		validator, err := utils.ConvertValidator(cdc, val, doc.InitialHeight)
		if err != nil {
			return err
		}

		vals[i] = validator
	}

	//return db.SaveValidatorsData(vals)
	// TODO
	_ = vals
	return nil
}

// saveValidatorDescription saves the description for the given validators
func saveValidatorDescription(doc *tmtypes.GenesisDoc, validators stakingtypes.Validators) error {
	descriptions := make([]types.ValidatorDescription, 0)
	for _, account := range validators {
		description, err := utils.ConvertValidatorDescription(
			account.OperatorAddress,
			account.Description,
			doc.InitialHeight,
		)
		if err != nil {
			return err
		}
		descriptions = append(descriptions, description)

		//err = db.SaveValidatorDescription(description)
		//if err != nil {
		//	return err
		//}
	}

	// TODO:
	_ = descriptions

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// saveDelegations stores the delegations data present inside the given genesis state
func saveDelegations(doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState) error {
	delegations := make([]types.Delegation, 0)
	for _, validator := range genState.Validators {
		tokens := validator.Tokens
		delegatorShares := validator.DelegatorShares

		for _, delegation := range findDelegations(genState.Delegations, validator.OperatorAddress) {
			delegationAmount := tokens.ToDec().Mul(delegation.Shares).Quo(delegatorShares).TruncateInt()
			delegations = append(delegations, types.NewDelegation(
				delegation.DelegatorAddress,
				validator.OperatorAddress,
				sdk.NewCoin(genState.Params.BondDenom, delegationAmount),
				doc.InitialHeight,
			))
		}
	}

	//if err := db.SaveDelegations(delegations); err != nil {
	//	return err
	//}
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

// saveUnbondingDelegations stores the unbonding delegations data present inside the given genesis state
func saveUnbondingDelegations(doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState) error {
	unbondingDelegations := make([]types.UnbondingDelegation, 0)
	for _, validator := range genState.Validators {
		valUD := findUnbondingDelegations(genState.UnbondingDelegations, validator.OperatorAddress)
		for _, ud := range valUD {
			for _, entry := range ud.Entries {
				unbondingDelegations = append(unbondingDelegations, types.NewUnbondingDelegation(
					ud.DelegatorAddress,
					validator.OperatorAddress,
					sdk.NewCoin(genState.Params.BondDenom, entry.InitialBalance),
					entry.CompletionTime,
					doc.InitialHeight,
				))
			}
		}
	}

	//return db.SaveUnbondingDelegations(unbondingDelegations)
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

// saveRedelegations stores the redelegations data present inside the given genesis state
func saveRedelegations(doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState) error {
	redelegations := make([]types.Redelegation, 0)
	for _, redelegation := range genState.Redelegations {
		for _, entry := range redelegation.Entries {
			redelegations = append(redelegations, types.NewRedelegation(
				redelegation.DelegatorAddress,
				redelegation.ValidatorSrcAddress,
				redelegation.ValidatorDstAddress,
				sdk.NewCoin(genState.Params.BondDenom, entry.InitialBalance),
				entry.CompletionTime,
				doc.InitialHeight,
			))
		}
	}

	// db.SaveRedelegations(redelegations)
	// TODO:
	_ = redelegations
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// saveValidatorsCommissions save the initial commission for each validator
func saveValidatorsCommissions(height int64, validators stakingtypes.Validators) error {
	validatorCommissions := make([]types.ValidatorCommission, len(validators))
	for i, account := range validators {
		validatorCommissions[i] = types.NewValidatorCommission(
			account.OperatorAddress,
			&account.Commission.Rate,
			&account.MinSelfDelegation,
			height,
		)
		//err := db.SaveValidatorCommission()
		//if err != nil {
		//	return err
		//}
	}

	return nil
}
