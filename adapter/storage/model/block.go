package model

import (
	"log"
	"time"
)

const (
	StatusProcessing Status = 1 + iota
	StatusProcessed
	StatusError
)

type (
	Status uint8

	Block struct {
		Processed    *time.Time `bson:"processed"`
		Created      time.Time
		ErrorMessage string `bson:"error_message"`
		Height       int64  `bson:"height"`
		Status       Status
	}
)

func (s Status) ToString() string {
	switch s {
	case StatusProcessing:
		return "processing"
	case StatusProcessed:
		return "processed"
	case StatusError:
		return "error"
	}

	log.Fatalf("uncnown status:%v", s)
	return ""
}

func (s Status) IsProcessing() bool { return s == StatusProcessing }
func (s Status) IsProcessed() bool  { return s == StatusProcessed }
func (s Status) IsError() bool      { return s == StatusError }
