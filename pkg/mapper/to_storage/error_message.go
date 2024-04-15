package tostorage

import (
	"time"

	"github.com/bro-n-bro/spacebox-crawler/v2/adapter/storage/model"
)

func (ts *ToStorage) NewErrorMessage(height int64, errMsg string) model.Message {
	return model.Message{
		Height:       height,
		Created:      time.Now(),
		ErrorMessage: errMsg,
	}
}
