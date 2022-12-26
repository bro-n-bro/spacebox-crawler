package storage

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"bro-n-bro-osmosis/adapter/storage/model"
)

func (s *Storage) HasBlock(ctx context.Context, height int64) (bool, error) {
	block := model.Block{}
	err := s.collection.FindOne(ctx, bson.D{{"height", height}}).Decode(&block)
	if err == nil {
		return true, nil

	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	return false, err
}

func (s *Storage) GetBlockStatus(ctx context.Context, height int64) (model.Status, error) {
	block := model.Block{}
	err := s.collection.FindOne(ctx, bson.D{{"height", height}}).Decode(&block)
	if err != nil {
		return 0, err

	}
	return block.Status, nil
}

func (s *Storage) CreateBlock(ctx context.Context, block *model.Block) error {
	_, err := s.collection.InsertOne(ctx, block)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) SetProcessedStatus(ctx context.Context, height int64) error {
	processed := time.Now()
	filter := bson.D{{"height", height}}
	update := bson.D{
		{"$set", bson.D{
			{"status", model.StatusProcessed},
			{"processed", &processed},
		}}}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) SetErrorStatus(ctx context.Context, height int64) error {
	filter := bson.D{{"height", height}}
	update := bson.D{
		{"$set", bson.D{
			{"status", model.StatusError},
		}}}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateStatus(ctx context.Context, height int64, status model.Status) error {
	filter := bson.D{{"height", height}}
	update := bson.D{{"$set", bson.D{{"status", status}}}}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetErrorBlockHeights(ctx context.Context) ([]int64, error) {
	cursor, err := s.collection.Find(ctx, bson.D{{"status", model.StatusError}})
	if err != nil {
		return nil, err
	}

	blocks := make([]model.Block, 0)
	if err = cursor.All(ctx, &blocks); err != nil {
		return nil, err
	}

	res := make([]int64, len(blocks))
	for i, block := range blocks {
		res[i] = block.Height
	}

	return res, nil
}
func (s *Storage) setErrorStatusForProcessing(ctx context.Context) error {
	filter := bson.D{{"status", model.StatusProcessing}}
	update := bson.D{{"$set", bson.D{{"status", model.StatusError}}}}
	_, err := s.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
