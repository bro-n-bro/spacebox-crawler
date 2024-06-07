package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/bro-n-bro/spacebox-crawler/v2/adapter/storage/model"
)

func (s *Storage) InsertErrorMessage(ctx context.Context, message model.Message) error {
	if _, err := s.messagesCollection.InsertOne(ctx, message); err != nil {
		return err
	}

	return nil
}

func (s *Storage) CountErrorMessages(ctx context.Context) (int64, error) {
	return s.messagesCollection.CountDocuments(ctx, bson.D{})
}
