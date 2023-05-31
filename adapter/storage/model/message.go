package model

import "time"

type Message struct {
	Created      time.Time
	ErrorMessage string `bson:"error_message"`
	Height       int64  `bson:"height"`
}
