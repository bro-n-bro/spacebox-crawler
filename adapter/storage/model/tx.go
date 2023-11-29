package model

import "time"

type Tx struct {
	Created      time.Time
	ErrorMessage string `bson:"error_message"`
	Hash         string `bson:"hash"`
	Height       int64  `bson:"height"`
}
