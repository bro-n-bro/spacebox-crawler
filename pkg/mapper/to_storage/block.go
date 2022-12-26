package to_storage

import (
	"time"

	"bro-n-bro-osmosis/adapter/storage/model"
)

func (ts ToStorage) NewBlock(height int64) *model.Block {
	return &model.Block{
		Height:  height,
		Created: time.Now(),
		Status:  model.StatusProcessing,
	}
}
