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

type DriverResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type NearestDriver struct {
	DriverID primitive.ObjectID `json:"id,omitempty"`
	Distance float64            `json:"distance,omitempty"`
}

type DriverUpsertRequest struct {
	Id       string          `json:"id,omitempty"`
	Location GeoJSONLocation `json:"location,omitempty" validate:"required"`
}
