package mongo

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/gt-tallinn/request-explorer/handlers/add"
	"context"
)

type Storage struct {
	collection *mongo.Collection
}

func New(collection *mongo.Collection) *Storage {
	return &Storage{
		collection: collection,
	}
}

func (s *Storage) Write(ctx context.Context, req *add.Request) error {
	_, err := s.collection.InsertOne(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

