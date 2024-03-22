package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Driver struct {
	ID       primitive.ObjectID `json:"id,omitempty"`
	Location GeoJSONLocation    `json:"location,omitempty" validate:"required"`
}

type GeoJSONLocation struct {
	Type        string    `json:"type,omitempty" validate:"required"`
	Coordinates []float64 `json:"coordinates,omitempty"`
}
