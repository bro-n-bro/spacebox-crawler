package storage

import (
	"context"
	"math/rand"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type block struct {
	Height int64 `bson:"height"`
}

func fillData(b *testing.B, maxRows int64, collectionA, collectionB *mongo.Collection) {
	b.Helper()

	mod := mongo.IndexModel{
		Keys: bson.M{"height": 1}, // index in ascending order or -1 for descending order
	}

	if _, err := collectionB.Indexes().CreateOne(context.Background(), mod); err != nil {
		b.Fatal(err)
	}

	blocks := make([]interface{}, maxRows)

	for i := int64(0); i < maxRows; i++ {
		blocks[i] = &block{i}
	}

	_, err := collectionA.InsertMany(context.Background(), blocks)
	if err != nil {
		b.Fatal(err)
	}

	_, err = collectionB.InsertMany(context.Background(), blocks)
	if err != nil {
		b.Fatal(err)
	}
}

// BenchmarkStorage/read_without_index
// BenchmarkStorage/read_without_index-8         	       1	2215740916 ns/op
// BenchmarkStorage/read_with_index
// BenchmarkStorage/read_with_index-8            	     166	   7426520 ns/op
// BenchmarkStorage/write_without_index
// BenchmarkStorage/write_without_index-8        	    1784	    749931 ns/op
// BenchmarkStorage/write_with_index
// BenchmarkStorage/write_with_index-8           	    1497	    732784 ns/op
func BenchmarkStorage(b *testing.B) {
	opts := []*options.ClientOptions{
		options.Client().ApplyURI("mongodb://localhost:27018/spacebox"),
		options.Client().SetMaxPoolSize(8),
		options.Client().SetMaxConnecting(8),
		options.Client().SetAuth(options.Credential{
			Username: "spacebox_user",
			Password: "spacebox_password",
		}),
	}

	client, err := mongo.Connect(context.Background(), opts...)
	if err != nil {
		b.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	if err := client.Ping(context.Background(), nil); err != nil {
		b.Fatal(err)
	}

	collectionA := client.Database("spacebox").Collection("blocksA")
	collectionB := client.Database("spacebox").Collection("blocksB")

	maxRows := int64(5000000)

	toSearchIndexes := 10
	indexes := make([]int64, 0, toSearchIndexes)

	for i := 0; i < toSearchIndexes; i++ {
		indexes = append(indexes, rand.Int63n(maxRows))
	}

	toWriteHeights := 100
	writeHeights := make([]int64, 0, toWriteHeights)
	for i := 0; i < toWriteHeights; i++ {
		writeHeights = append(writeHeights, maxRows+int64(i+1))
	}

	fillData(b, maxRows, collectionA, collectionB)

	b.Run("read without index", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, height := range indexes {
				if err := collectionA.
					FindOne(context.Background(), bson.D{{Key: "height", Value: height}}).
					Decode(&block{}); err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("read with index", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, height := range indexes {
				if err := collectionB.
					FindOne(context.Background(), bson.D{{Key: "height", Value: height}}).
					Decode(&block{}); err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("write without index", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h := writeHeights[rand.Intn(len(writeHeights)-1)]
			if _, err := collectionA.InsertOne(context.Background(), &block{h}); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("write with index", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h := writeHeights[rand.Intn(len(writeHeights)-1)]
			if _, err := collectionB.InsertOne(context.Background(), &block{h}); err != nil {
				b.Fatal(err)
			}
		}
	})
}
