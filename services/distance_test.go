package services

import (
	"math"
	"testing"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
		delta    float64 // 許容誤差
	}{
		{
			name:     "同一地点",
			lat1:     35.6762,
			lon1:     139.6503,
			lat2:     35.6762,
			lon2:     139.6503,
			expected: 0,
			delta:    0.001,
		},
		{
			name:     "東京-大阪",
			lat1:     35.6762,
			lon1:     139.6503,
			lat2:     34.6937,
			lon2:     135.5023,
			expected: 395,
			delta:    5, // 5km誤差許容
		},
		{
			name:     "北半球-南半球（東京-シドニー）",
			lat1:     35.6762,
			lon1:     139.6503,
			lat2:     -33.8688,
			lon2:     151.2093,
			expected: 7823,
			delta:    50, // 50km誤差許容
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.Abs(got-tt.expected) > tt.delta {
				t.Errorf("HaversineDistance() = %v, want %v (±%v)", got, tt.expected, tt.delta)
			}
		})
	}
}
