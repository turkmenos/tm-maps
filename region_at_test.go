package tmmaps

import (
	"encoding/json"
	"errors"
	"math"
	"testing"
)

func TestRegionAt(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		wantSlug  string
	}{
		{"Ashgabat is covered by Ahal in this dataset", 37.960077, 58.326063, "turkmenistan-ahal"},
		{"Balkanabat", 39.51075, 54.36713, "turkmenistan-balkan"},
		{"Dashoguz", 41.83625, 59.96661, "turkmenistan-dasoguz"},
		{"Turkmenabat", 39.07328, 63.57861, "turkmenistan-lebap"},
		{"Mary", 37.59378, 61.83031, "turkmenistan-mary"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			region, err := RegionAt(tt.latitude, tt.longitude)
			if err != nil {
				t.Fatal(err)
			}
			if region.Slug != tt.wantSlug {
				t.Fatalf("RegionAt() slug = %q, want %q", region.Slug, tt.wantSlug)
			}
		})
	}
}

func TestRegionAtOutside(t *testing.T) {
	_, err := RegionAt(0, 0)
	if !errors.Is(err, ErrRegionNotFound) {
		t.Fatalf("RegionAt() error = %v, want ErrRegionNotFound", err)
	}
}

func TestRegionAtInvalidCoordinates(t *testing.T) {
	for _, coordinate := range [][2]float64{{91, 0}, {-91, 0}, {0, 181}, {0, -181}, {math.NaN(), 0}} {
		_, err := RegionAt(coordinate[0], coordinate[1])
		if !errors.Is(err, ErrInvalidCoordinate) {
			t.Errorf("RegionAt(%v, %v) error = %v, want ErrInvalidCoordinate", coordinate[0], coordinate[1], err)
		}
	}
}

func TestRegionAtBoundaryPoint(t *testing.T) {
	data, err := Welaýat("ahal")
	if err != nil {
		t.Fatal(err)
	}
	var collection struct {
		Features []struct {
			Geometry struct {
				Coordinates [][][]float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}
	if err := json.Unmarshal(data, &collection); err != nil {
		t.Fatal(err)
	}
	point := collection.Features[0].Geometry.Coordinates[0][0]
	region, err := RegionAt(point[1], point[0])
	if err != nil {
		t.Fatal(err)
	}
	if region.Slug == "" {
		t.Fatal("RegionAt() returned an empty region on a boundary point")
	}
}

func TestGeometryContainsMultiPolygon(t *testing.T) {
	coordinates := [][][][]float64{{{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}}}, {{{10, 10}, {12, 10}, {12, 12}, {10, 12}, {10, 10}}}}
	data, err := json.Marshal(coordinates)
	if err != nil {
		t.Fatal(err)
	}
	inside, err := geometryContains(boundaryGeometry{Type: "MultiPolygon", Coordinates: data}, [2]float64{11, 11})
	if err != nil {
		t.Fatal(err)
	}
	if !inside {
		t.Fatal("geometryContains() did not match a MultiPolygon point")
	}
}
