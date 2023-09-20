package utils

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-crawler/types"
	"github.com/bro-n-bro/spacebox/broker/model"
)

type (
	AccountCache[K, V comparable] interface {
		UpdateCacheValue(K, V) bool
	}

	AccountPublisher interface {
		PublishAccount(ctx context.Context, account model.Account) error
	}
)

// GetGenesisAccounts parses the given appState and returns the genesis accounts
func GetGenesisAccounts(appState map[string]json.RawMessage, cdc codec.Codec) ([]types.Account, error) {
	var authState authtypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[authtypes.ModuleName], &authState); err != nil {
		return nil, err
	}

	// Store the accounts
	accounts := make([]types.Account, 0)

	for _, account := range authState.Accounts {
		var accountI authtypes.AccountI
		if err := cdc.UnpackAny(account, &accountI); err != nil {
			return nil, err
		}

		accounts = append(accounts, types.NewAccount(accountI.GetAddress().String(), account.TypeUrl, 0))
	}

	return accounts, nil
}

// FilterNonAccountAddresses returns a slice containing only account addresses.
func FilterNonAccountAddresses(addresses []string) []string {
	// Filter using only the account addresses as the MessageAddressesParser might return also validator addresses
	accountAddresses := make([]string, 0)

	for _, address := range addresses {
		if _, err := sdk.AccAddressFromBech32(address); err == nil { // needs correct addresses only
			accountAddresses = append(accountAddresses, address)
		}
	}

	return accountAddresses
}

// GetAndPublishAccount retrieves the account from the given address and publishes it.
func GetAndPublishAccount(
	ctx context.Context,
	address string,
	height int64,
	cache AccountCache[string, int64],
	publisher AccountPublisher,
	client authtypes.QueryClient,
) error {

	// allow only for account addresses
	if !strings.HasPrefix(address, sdk.GetConfig().GetBech32AccountAddrPrefix()) ||
		strings.HasPrefix(address, sdk.GetConfig().GetBech32ValidatorAddrPrefix()) ||
		strings.HasPrefix(address, sdk.GetConfig().GetBech32ConsensusAddrPrefix()) {

		return nil
	}

	if !cache.UpdateCacheValue(address, height) {
		return nil
	}

	respPb, err := client.Account(ctx, &authtypes.QueryAccountRequest{Address: address})
	if err != nil {
		return errors.Wrap(err, "fail to get account")
	}

	return publisher.PublishAccount(ctx, model.Account{
		Address: address,
		Type:    respPb.Account.TypeUrl,
		Height:  height,
	})
}
