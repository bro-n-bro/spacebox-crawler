package model

import "time"

type Message struct {
	Height       int64 `bson:"height"`
	Created      time.Time
	ErrorMessage string `bson:"error_message"`
}
