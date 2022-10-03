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

var conCity = CityDAO{}

func init() {
	conCity.Server = "mongodb://localhost:27017/"
	conCity.Database = "CityDB"
	conCity.Collection = "City"

	conCity.Connect()
}

type CityDAO struct {
	Server     string
	Database   string
	Collection string
}

var CollectionCity *mongo.Collection
var ctxCity = context.TODO()

func (e *CityDAO) Connect() {
	clientOptions := options.Client().ApplyURI(e.Server)
	client, err := mongo.Connect(ctxCity, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctxCity, nil)
	if err != nil {
		log.Fatal(err)
	}

	CollectionCity = client.Database(e.Database).Collection(e.Collection)
}

func (e *CityDAO) Insert(city modelPojo.City) error {
	_, err := CollectionCity.InsertOne(ctxCity, city)

	if err != nil {
		return errors.New("unable to create new record")
	}

	return nil
}

func (e *CityDAO) FindByCityName(cityName string) ([]*modelPojo.City, error) {
	var city []*modelPojo.City

	cur, err := CollectionCity.Find(ctxCity, bson.D{primitive.E{Key: "city_name", Value: cityName}})

	if err != nil {
		return city, errors.New("unable to query db")
	}

	for cur.Next(ctxCity) {
		var e modelPojo.City

		err := cur.Decode(&e)

		if err != nil {
			return city, err
		}

		city = append(city, &e)
	}

	if err := cur.Err(); err != nil {
		return city, err
	}

	cur.Close(ctxCity)

	if len(city) == 0 {
		return city, mongo.ErrNoDocuments
	}

	return city, nil
}

func (e *CityDAO) DeleteCity(cityName string) error {
	filter := bson.D{primitive.E{Key: "city_name", Value: cityName}}

	res, err := CollectionCity.DeleteOne(ctxCity, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no city deleted")
	}

	return nil
}

func (epd *CityDAO) UpdateCity(cityName string, city modelPojo.City) error {
	filter := bson.D{primitive.E{Key: "city_name", Value: cityName}}

	update := bson.D{primitive.E{Key: "$set", Value: city}}

	e := &modelPojo.City{}
	return CollectionCity.FindOneAndUpdate(ctxCity, filter, update).Decode(e)
}
