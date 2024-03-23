package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type RiderRequest struct {
	ID              primitive.ObjectID `json:"id,omitempty"`
	Type            string             `json:"type,omitempty" validate:"required"`
	UserCoordinates []float64          `json:"coordinates,omitempty" validate:"required"`
	UserRadius      float64            `json:"radius,omitempty" validate:"required"`
}

type RiderResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
