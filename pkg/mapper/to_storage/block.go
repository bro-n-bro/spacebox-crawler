package tostorage

import (
	"time"

	"github.com/bro-n-bro/spacebox-crawler/adapter/storage/model"
)

func (ts ToStorage) NewBlock(height int64) *model.Block {
	return &model.Block{
		Height:  height,
		Created: time.Now(),
		Status:  model.StatusProcessing,
	}
}
