package distribution

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bro-n-bro/spacebox-crawler/types"
)

func TestCoinFromBytes(t *testing.T) {
	testCases := []struct {
		value string
		want  types.Coins
	}{
		{
			value: "12345uatom",
			want: types.Coins{
				types.Coin{
					Denom:  "uatom",
					Amount: 12345,
				},
			},
		},
		{
			value: "123456789uatom",
			want: types.Coins{
				types.Coin{
					Denom:  "uatom",
					Amount: 123456789,
				},
			},
		},
		{
			value: "1uatom",
			want: types.Coins{
				types.Coin{
					Denom:  "uatom",
					Amount: 1,
				},
			},
		},
		{
			value: "1ibc/B05539B66B72E2739B986B86391E5D08F12B8D5D2C2A7F8F8CF9ADF674DFA231,4146906uatom",
			want: types.Coins{
				types.Coin{
					Denom:  "ibc/B05539B66B72E2739B986B86391E5D08F12B8D5D2C2A7F8F8CF9ADF674DFA231",
					Amount: 1,
				},
				types.Coin{
					Denom:  "uatom",
					Amount: 4146906,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.value, func(t *testing.T) {
			got, err := coinsFromAttribute(tt.value)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
