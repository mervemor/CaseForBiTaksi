package helpers

import (
	"testing"
)

func TestHaversine(t *testing.T) {
	tests := []struct {
		name           string
		latitude1      float64
		longitude1     float64
		latitude2      float64
		longitude2     float64
		expectedResult float64
	}{
		{
			name:           "Distance between two same points",
			latitude1:      40.94289771,
			longitude1:     29.0390297,
			latitude2:      40.94289771,
			longitude2:     29.0390297,
			expectedResult: 0,
		},
		{
			name:           "Distance between two different points",
			latitude1:      40.58889619,
			longitude1:     29.4638355,
			latitude2:      40.581087,
			longitude2:     29.46146,
			expectedResult: 0.89,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Haversine(tt.latitude1, tt.longitude1, tt.latitude2, tt.longitude2)
			if got != tt.expectedResult {
				t.Errorf("Haversine() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}
