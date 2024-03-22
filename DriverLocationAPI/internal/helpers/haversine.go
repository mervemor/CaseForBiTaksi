package helpers

import (
	"math"
)

const (
	EarthRadius = 6371
)

func Haversine(latitude1, longitude1, latitude2, longitude2 float64) float64 {
	lat1 := degreesToRadians(latitude1)
	lon1 := degreesToRadians(longitude1)
	lat2 := degreesToRadians(latitude2)
	lon2 := degreesToRadians(longitude2)

	dLat := lat2 - lat1
	dLon := lon2 - lon1

	a := math.Pow(math.Sin(dLat/2), 2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := EarthRadius * c
	return math.Round(distance*100) / 100
}

func degreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}
