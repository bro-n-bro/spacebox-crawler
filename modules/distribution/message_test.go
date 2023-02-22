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
	}

	for _, tt := range testCases {
		t.Run(tt.value, func(t *testing.T) {
			got, err := coinFromBytes([]byte(tt.value))
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
