package storage

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bro-n-bro/spacebox-crawler/adapter/storage/model"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

func (s *Storage) GetBlockByHeight(ctx context.Context, height int64) (*model.Block, error) {
	var (
		err   error
		block model.Block
	)

	if err = s.collection.
		FindOne(ctx, bson.D{{Key: "height", Value: height}}).
		Decode(&block); err == nil { // record exist
		return &block, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, types.ErrBlockNotFound
	}

	return nil, err
}

func (s *Storage) CreateBlock(ctx context.Context, block *model.Block) error {
	if _, err := s.collection.InsertOne(ctx, block); err != nil {
		return err
	}

	return nil
}

func (s *Storage) SetProcessedStatus(ctx context.Context, height int64) error {
	processed := time.Now()
	filter := bson.D{{Key: "height", Value: height}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: model.StatusProcessed},
			{Key: "processed", Value: &processed},
			{Key: "error_message", Value: ""},
		}}}
	if _, err := s.collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (s *Storage) SetErrorStatus(ctx context.Context, height int64, msg string) error {
	filter := bson.D{{Key: "height", Value: height}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: model.StatusError},
			{Key: "error_message", Value: msg},
		}}}
	if _, err := s.collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateStatus(ctx context.Context, height int64, status model.Status) error {
	filter := bson.D{{Key: "height", Value: height}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: status}}}}
	if _, err := s.collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetErrorBlockHeights(ctx context.Context) ([]int64, error) {
	cursor, err := s.collection.Find(ctx, bson.D{{Key: "status", Value: model.StatusError}})
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

func (s *Storage) GetAllBlocks(ctx context.Context) ([]*model.Block, error) {
	cursor, err := s.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			s.log.Error().Err(err).Msg("close cursor error")
		}
	}()

	blocks := make([]*model.Block, 0)
	if err = cursor.All(ctx, &blocks); err != nil {
		return nil, err
	}

	return blocks, err
}

func (s *Storage) setErrorStatusForProcessing(ctx context.Context) error {
	filter := bson.D{{Key: "status", Value: model.StatusProcessing}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "status", Value: model.StatusError},
		{Key: "error_message", Value: "dont have time to process"},
	}}}

	if _, err := s.collection.UpdateMany(ctx, filter, update); err != nil {
		return err
	}

	return nil
}
