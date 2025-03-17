package persistence

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type InMongoDB[T any] struct {
	collectionName string
	connection     string
	client         *mongo.Client
	database       *mongo.Database
}

func NewInMongoDB[T any](collectionName, connection string) (*InMongoDB[T], error) {

	cs, err := connstring.ParseAndValidate(connection)
	if err != nil {
		return nil, err
	}

	databaseName := cs.Database
	if databaseName == "" {
		databaseName = "itemstest"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	if err != nil {
		return nil, err
	}

	// ensure unique "id" index for items :D
	database := client.Database(databaseName)
	_, err = database.Collection(collectionName).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{"id": 1},
	})
	if err != nil {
		return nil, err
	}

	return &InMongoDB[T]{
		collectionName: collectionName,
		connection:     connection,
		client:         client, // might not be needed
		database:       database,
	}, nil
}

func (f *InMongoDB[T]) List(ctx context.Context) ([]*ItemWithId[T], error) {

	cur, err := f.database.Collection(f.collectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	result := []*ItemWithId[T]{}

	for cur.Next(context.Background()) {
		item := &ItemWithId[T]{}
		err := cur.Decode(item)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (f *InMongoDB[T]) Put(ctx context.Context, item *ItemWithId[T]) error {
	filter := bson.M{
		"_id": item.Id,
	}
	if item.Version > 0 {
		filter["version"] = item.Version

	}
	update := bson.M{
		"$set": ItemWithId[T]{
			Id:      item.Id,
			Item:    item.Item,
			Version: item.Version + 1,
		},
	}

	result, err := f.database.Collection(f.collectionName).UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

	if mongo.IsDuplicateKeyError(err) {
		return ErrVersionGone
	}

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return ErrVersionGone
	}

	return nil
}

func (f *InMongoDB[T]) Get(ctx context.Context, id string) (*ItemWithId[T], error) {
	result := &ItemWithId[T]{}
	err := f.database.Collection(f.collectionName).FindOne(ctx, bson.M{"id": id}).Decode(result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return result, err
}

func (f *InMongoDB[T]) Delete(ctx context.Context, id string) error {
	_, err := f.database.Collection(f.collectionName).DeleteOne(ctx, bson.M{"id": id})
	return err
}
