package modelPojo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type City struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CityName string             `bson:"city_name,omitempty" json:"city_name,omitempty"`
	CityCode string             `bson:"city_code, omitempty" json:"city_code,omitempty"`
	State    string             `bson:"state, omitempty" json:"state,omitempty"`
	Country  string             `bson:"country, omitempty" json:"country,omitempty"`
	PinCode  int64              `bson:"pinCode, omitempty" json:"pinCode,omitempty"`
}

type Classification struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ServiceType string             `bson:"service_type,omitempty" json:"service_type,omitempty"`
}

type Service struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name           string             `bson:"name,omitempty" json:"name,omitempty"`
	Address        string             `bson:"address, omitempty" json:"address,omitempty"`
	Latitude       float64            `bson:"latitude, omitempty" json:"latitude,omitempty"`
	Longitude      float64            `bson:"longitude, omitempty" json:"longitude,omitempty"`
	Website        string             `bson:"website, omitempty" json:"website,omitempty"`
	ContactNumber  int64              `bson:"contact_number, omitempty" json:"contact_number,omitempty"`
	City           *City              `bson:"city, omitempty" json:"city,omitempty"`
	Verified       bool               `bson:"verified, omitempty" json:"verified,omitempty"`
	Classification *Classification    `bson:"classification, omitempty" json:"classification,omitempty"`
}

type Search struct {
	CityName    string `bson:"city_name,omitempty" json:"city_name,omitempty"`
	ServiceType string `bson:"service_type,omitempty" json:"service_type,omitempty"`
}
