package domain

type UserRequest struct {
	ID              string    `json:"id,omitempty"`
	Type            string    `json:"type,omitempty" validate:"required"`
	UserCoordinates []float64 `json:"coordinates,omitempty"`
	UserRadius      float64   `json:"radius,omitempty"`
}
