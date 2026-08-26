package tmmaps

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name      string
		slug      string
		latitude  float64
		longitude float64
		want      bool
	}{
		{"inside requested region", "ahal", 37.960077, 58.326063, true},
		{"inside a different region", "mary", 37.960077, 58.326063, false},
		{"outside Turkmenistan", "ahal", 0, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Contains(tt.slug, tt.latitude, tt.longitude)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("Contains(%q, %v, %v) = %v, want %v", tt.slug, tt.latitude, tt.longitude, got, tt.want)
			}
		})
	}
}

func TestContainsUnknownRegion(t *testing.T) {
	inside, err := Contains("unknown", 37.960077, 58.326063)
	if inside || !errors.Is(err, ErrUnknownRegion) {
		t.Fatalf("Contains() = %v, %v; want false, ErrUnknownRegion", inside, err)
	}
}

func TestContainsBoundaryPoint(t *testing.T) {
	data, err := Welaýat("ahal")
	if err != nil {
		t.Fatal(err)
	}
	var collection boundaryCollection
	if err := json.Unmarshal(data, &collection); err != nil {
		t.Fatal(err)
	}
	var polygon [][][]float64
	if err := json.Unmarshal(collection.Features[0].Geometry.Coordinates, &polygon); err != nil {
		t.Fatal(err)
	}
	point := polygon[0][0]
	inside, err := Contains("ahal", point[1], point[0])
	if err != nil {
		t.Fatal(err)
	}
	if !inside {
		t.Fatal("Contains() returned false for a boundary point")
	}
}
