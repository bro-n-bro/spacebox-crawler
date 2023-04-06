package utils

import (
	"regexp"
	"strconv"
	"strings"

	"cosmossdk.io/errors"

	"github.com/bro-n-bro/spacebox-crawler/types"
)

var (
	withdrawDelegationRewardRegex = regexp.MustCompile(`^(\-?[0-9]+(\.[0-9]+)?)([0-9a-zA-Z/]+)$`)
)

// ParseCoinsFromString converts string to coin type
func ParseCoinsFromString(value string) (types.Coins, error) {
	rows := strings.Split(value, ",")
	res := make(types.Coins, len(rows))
	for i, row := range rows {
		bits := withdrawDelegationRewardRegex.FindStringSubmatch(row)
		if len(bits) < 4 {
			continue
		}
		amount, err := strconv.ParseFloat(bits[1], 64)
		if err != nil {
			return types.Coins{}, errors.Wrap(err, "failed to parse float")
		}
		res[i] = types.Coin{
			Denom:  bits[3],
			Amount: amount,
		}
	}
	return res, nil
}
