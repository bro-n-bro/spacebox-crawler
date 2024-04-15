package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/bro-n-bro/spacebox-crawler/v2/adapter/storage/model"
)

func (s *Storage) InsertErrorTx(ctx context.Context, tx model.Tx) error {
	if _, err := s.txCollection.InsertOne(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (s *Storage) CountErrorTxs(ctx context.Context) (int64, error) {
	return s.txCollection.CountDocuments(ctx, bson.D{})
}
