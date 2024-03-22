package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserRequest struct {
	ID              string    `json:"id,omitempty"`
	Type            string    `json:"type,omitempty" validate:"required"`
	UserCoordinates []float64 `json:"coordinates,omitempty"`
	UserRadius      float64   `json:"radius,omitempty"`
}

type DistanceBetweenDriverAndUser struct {
	DriverID primitive.ObjectID `json:"id,omitempty"`
	Distance float64            `json:"distance,omitempty"`
}
