package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *Client, collName string) *Repository {
	collection := db.Database.Collection("urls")
	return &Repository{collection}
}

func (m *Repository) FindAll(ctx context.Context, filter interface{}, v []interface{}) ([]interface{}, error) {
	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (m *Repository) CreateDocument(ctx context.Context, d interface{}) (primitive.ObjectID, error) {
	if d == nil {
		return primitive.NilObjectID, fmt.Errorf("document cannot be nil")
	}

	res, err := m.collection.InsertOne(ctx, d)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if res == nil {
		return primitive.NilObjectID, nil
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (m *Repository) Count(ctx context.Context, f interface{}) (int64, error) {
	return m.collection.CountDocuments(ctx, f)
}

func (m *Repository) FindOne(ctx context.Context, f interface{}, v interface{}) error {
	return m.collection.FindOne(ctx, f).Decode(v)
}

func (m *Repository) Delete(ctx context.Context, f interface{}) error {
	_, err := m.collection.DeleteMany(ctx, f)
	return err
}

func (m *Repository) Update(ctx context.Context, oid primitive.ObjectID, d interface{}) error {
	res := m.collection.FindOneAndUpdate(
		ctx,
		primitive.M{"_id": oid},
		primitive.M{"$set": d}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	if res.Err() != nil {
		return fmt.Errorf("failed to perform find one and update: %v", res.Err())
	}

	if err := res.Decode(d); err != nil {
		return fmt.Errorf("failed to decode: %v", err)
	}
	return nil
}
