package dao

import (
	"context"
	"errors"
	"log"

	"city-search-project/modelPojo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var con = CategoryDAO{}

func init() {
	con.Server = "mongodb://localhost:27017/"
	con.Database = "CityDB"
	con.Collection = "Category"

	con.Connect()
}

type CategoryDAO struct {
	Server     string
	Database   string
	Collection string
}

var Collection *mongo.Collection
var ctx = context.TODO()

func (e *CategoryDAO) Connect() {
	clientOptions := options.Client().ApplyURI(e.Server)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	Collection = client.Database(e.Database).Collection(e.Collection)
}

func (e *CategoryDAO) Insert(category modelPojo.Classification) error {
	_, err := Collection.InsertOne(ctx, category)

	if err != nil {
		return errors.New("unable to create new record")
	}

	return nil
}

func (e *CategoryDAO) FindByCategory(serviceType string) ([]*modelPojo.Classification, error) {
	var category []*modelPojo.Classification

	cur, err := Collection.Find(ctx, bson.D{primitive.E{Key: "service_type", Value: serviceType}})

	if err != nil {
		return category, errors.New("unable to query db")
	}

	for cur.Next(ctx) {
		var e modelPojo.Classification

		err := cur.Decode(&e)

		if err != nil {
			return category, err
		}

		category = append(category, &e)
	}

	if err := cur.Err(); err != nil {
		return category, err
	}

	cur.Close(ctx)

	if len(category) == 0 {
		return category, mongo.ErrNoDocuments
	}

	return category, nil
}

func (e *CategoryDAO) DeleteCategory(serviceType string) error {
	filter := bson.D{primitive.E{Key: "service_type", Value: serviceType}}

	res, err := Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no category deleted")
	}

	return nil
}

func (epd *CategoryDAO) UpdateCategory(serviceType string, category modelPojo.Classification) error {
	filter := bson.D{primitive.E{Key: "service_type", Value: serviceType}}

	update := bson.D{primitive.E{Key: "$set", Value: category}}

	e := &modelPojo.Classification{}
	return Collection.FindOneAndUpdate(ctx, filter, update).Decode(e)
}
