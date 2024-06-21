package persistence

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type InMongoDB struct {
	connection string
	client     *mongo.Client
	database   *mongo.Database
}

func NewInMongoDB(connection string) (*InMongoDB, error) {

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
	_, err = database.Collection("items").Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{"id": 1},
	})
	if err != nil {
		return nil, err
	}

	return &InMongoDB{
		connection: connection,
		client:     client, // might not be needed
		database:   database,
	}, nil
}

func (f *InMongoDB) List(ctx context.Context) ([]*ItemWithId, error) {

	cur, err := f.database.Collection("items").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	result := []*ItemWithId{}

	for cur.Next(context.Background()) {
		item := &ItemWithId{}
		err := cur.Decode(item)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (f *InMongoDB) Put(ctx context.Context, item *ItemWithId) error {
	filter := bson.M{
		"_id": item.Id,
	}
	if item.Version > 0 {
		filter["version"] = item.Version

	}
	update := bson.M{
		"$set": ItemWithId{
			Id:      item.Id,
			Item:    item.Item,
			Version: item.Version + 1,
		},
	}

	result, err := f.database.Collection("items").UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

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

func (f *InMongoDB) Get(ctx context.Context, id string) (*ItemWithId, error) {
	result := &ItemWithId{}
	err := f.database.Collection("items").FindOne(ctx, bson.M{"id": id}).Decode(result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return result, err
}

func (f *InMongoDB) Delete(ctx context.Context, id string) error {
	_, err := f.database.Collection("items").DeleteOne(ctx, bson.M{"id": id})
	return err
}
