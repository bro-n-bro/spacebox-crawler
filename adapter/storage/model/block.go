package model

import "time"

type Status uint8

const (
	StatusProcessing Status = 1 + iota
	StatusProcessed
	StatusError
)

type Block struct {
	Processed    *time.Time `bson:"processed"`
	Created      time.Time
	ErrorMessage string `bson:"error_message"`
	Height       int64  `bson:"height"`
	Status       Status
}

func (s Status) IsProcessing() bool { return s == StatusProcessing }
func (s Status) IsProcessed() bool  { return s == StatusProcessed }
func (s Status) IsError() bool      { return s == StatusError }
