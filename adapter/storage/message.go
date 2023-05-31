package storage

import (
	"context"

	"github.com/bro-n-bro/spacebox-crawler/adapter/storage/model"
)

func (s *Storage) InsertErrorMessage(ctx context.Context, message model.Message) error {
	if _, err := s.messagesCollection.InsertOne(ctx, message); err != nil {
		return err
	}

	return nil
}

func (s *Storage) CountErrorMessage(ctx context.Context) (int64, error) {
	return s.messagesCollection.CountDocuments(ctx, nil)
}
