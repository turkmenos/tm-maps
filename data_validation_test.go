package tmmaps

import (
	"encoding/json"
	"math"
	"testing"
)

func TestRegionDataIntegrity(t *testing.T) {
	records, err := Regions()
	if err != nil {
		t.Fatal(err)
	}
	bySlug := make(map[string]Region, len(records))
	for _, record := range records {
		if record.Slug == "" {
			t.Fatal("region data contains an empty slug")
		}
		if _, exists := bySlug[record.Slug]; exists {
			t.Errorf("region data contains duplicate slug %q", record.Slug)
		}
		bySlug[record.Slug] = record
		if record.NameTM == "" || record.Type == "" {
			t.Errorf("region %q is missing name_tm or type", record.Slug)
		}
		if (record.Latitude == nil) != (record.Longitude == nil) {
			t.Errorf("region %q has an incomplete coordinate pair", record.Slug)
		}
		if record.Latitude != nil && !validDataCoordinate(*record.Latitude, *record.Longitude) {
			t.Errorf("region %q has invalid coordinates (%v, %v)", record.Slug, *record.Latitude, *record.Longitude)
		}
	}
	for _, record := range records {
		if record.ParentSlug == "" {
			if record.Type != "country" {
				t.Errorf("non-country region %q has no parent", record.Slug)
			}
			continue
		}
		if _, exists := bySlug[record.ParentSlug]; !exists {
			t.Errorf("region %q references unknown parent %q", record.Slug, record.ParentSlug)
		}
	}
}

func TestWelayatGeoJSONIntegrity(t *testing.T) {
	for _, name := range []string{"ahal", "balkan", "dasoguz", "lebap", "mary"} {
		t.Run(name, func(t *testing.T) {
			data, err := Welaýat(name)
			if err != nil {
				t.Fatal(err)
			}
			var collection boundaryCollection
			if err := json.Unmarshal(data, &collection); err != nil {
				t.Fatal(err)
			}
			if len(collection.Features) == 0 {
				t.Fatal("GeoJSON contains no features")
			}
			wantSlug := "turkmenistan-" + name
			for i, feature := range collection.Features {
				if feature.Properties.Slug != wantSlug {
					t.Errorf("feature %d slug = %q, want %q", i, feature.Properties.Slug, wantSlug)
				}
				validateBoundaryGeometry(t, feature.Geometry)
			}
		})
	}
}

func validateBoundaryGeometry(t *testing.T, geometry boundaryGeometry) {
	t.Helper()
	switch geometry.Type {
	case "Polygon":
		var polygon [][][]float64
		if err := json.Unmarshal(geometry.Coordinates, &polygon); err != nil {
			t.Fatal(err)
		}
		validatePolygonCoordinates(t, polygon)
	case "MultiPolygon":
		var polygons [][][][]float64
		if err := json.Unmarshal(geometry.Coordinates, &polygons); err != nil {
			t.Fatal(err)
		}
		if len(polygons) == 0 {
			t.Error("MultiPolygon contains no polygons")
		}
		for _, polygon := range polygons {
			validatePolygonCoordinates(t, polygon)
		}
	default:
		t.Errorf("unsupported geometry type %q", geometry.Type)
	}
}

func validatePolygonCoordinates(t *testing.T, polygon [][][]float64) {
	t.Helper()
	if len(polygon) == 0 {
		t.Error("Polygon contains no rings")
		return
	}
	for _, ring := range polygon {
		if len(ring) < 4 {
			t.Errorf("polygon ring contains %d points, want at least 4", len(ring))
			continue
		}
		first, last := ring[0], ring[len(ring)-1]
		if len(first) < 2 || len(last) < 2 || first[0] != last[0] || first[1] != last[1] {
			t.Error("polygon ring is not closed")
		}
		for _, coordinate := range ring {
			if len(coordinate) < 2 || !validDataCoordinate(coordinate[1], coordinate[0]) {
				t.Errorf("invalid GeoJSON coordinate %v", coordinate)
			}
		}
	}
}

func validDataCoordinate(latitude, longitude float64) bool {
	return !math.IsNaN(latitude) && !math.IsNaN(longitude) &&
		!math.IsInf(latitude, 0) && !math.IsInf(longitude, 0) &&
		latitude >= -90 && latitude <= 90 && longitude >= -180 && longitude <= 180
}
