package model

import (
	"time"
)

type RedelegationMessage struct {
	CompletionTime   time.Time `json:"completion_time"`
	Coin             Coin      `json:"coin"`
	DelegatorAddress string    `json:"delegator_address"`
	SrcValidator     string    `json:"src_validator"`
	DstValidator     string    `json:"dst_validator"`
	TxHash           string    `json:"tx_hash"`
	Height           int64     `json:"height"`
}
