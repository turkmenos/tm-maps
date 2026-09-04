package tmmaps

import (
	"errors"
	"testing"
)

func TestPOIsByCategory(t *testing.T) {
	cafes, err := POIs(POICategoryCafe)
	if err != nil {
		t.Fatal(err)
	}
	if len(cafes) != 376 {
		t.Fatalf("POIs(cafes) returned %d results, want 376", len(cafes))
	}
	if cafes[0].Category != POICategoryCafe {
		t.Fatalf("POIs(cafes) category = %q, want cafes", cafes[0].Category)
	}
	if cafes[0].Latitude == 0 || cafes[0].Longitude == 0 {
		t.Fatal("POIs(cafes) omitted coordinates")
	}
}

func TestPOIsAeroway(t *testing.T) {
	aeroways, err := POIs(POICategoryAeroway)
	if err != nil {
		t.Fatal(err)
	}
	if len(aeroways) != 19 {
		t.Fatalf("POIs(aeroway) returned %d results, want 19", len(aeroways))
	}
	if aeroways[0].Name != "Türkmenbaşy Halkara Howa Menzili" {
		t.Fatalf("POIs(aeroway) first name = %q, want Türkmenbaşy Halkara Howa Menzili", aeroways[0].Name)
	}
}

func TestAllPOIs(t *testing.T) {
	points, err := AllPOIs()
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 746 {
		t.Fatalf("AllPOIs() returned %d results, want 746", len(points))
	}
}

func TestSearchPOIs(t *testing.T) {
	results, err := SearchPOIs("dayanc", POICategoryHotel)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("SearchPOIs() returned %d results, want 1", len(results))
	}
	if results[0].Name != "Dayanc Hotel" {
		t.Fatalf("SearchPOIs() name = %q, want Dayanc Hotel", results[0].Name)
	}
}

func TestPOIsWithinRadius(t *testing.T) {
	results, err := POIsWithinRadius(37.885931, 58.3592016, 1, POICategoryHotel)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("POIsWithinRadius() returned no hotels")
	}
	if results[0].Name != "Dayanc Hotel" {
		t.Fatalf("POIsWithinRadius() first name = %q, want Dayanc Hotel", results[0].Name)
	}
	for i := 1; i < len(results); i++ {
		if results[i].DistanceKM < results[i-1].DistanceKM {
			t.Fatalf("POIsWithinRadius() results are not ordered at indexes %d and %d", i-1, i)
		}
	}
}

func TestPOIInvalidInputs(t *testing.T) {
	if _, err := POIs("parks"); !errors.Is(err, ErrUnknownPOICategory) {
		t.Fatalf("POIs() error = %v, want ErrUnknownPOICategory", err)
	}
	if _, err := SearchPOIs("x", "parks"); !errors.Is(err, ErrUnknownPOICategory) {
		t.Fatalf("SearchPOIs() error = %v, want ErrUnknownPOICategory", err)
	}
	if _, err := POIsWithinRadius(91, 0, 1); !errors.Is(err, ErrInvalidCoordinate) {
		t.Fatalf("POIsWithinRadius() coordinate error = %v, want ErrInvalidCoordinate", err)
	}
	if _, err := POIsWithinRadius(0, 0, 0); !errors.Is(err, ErrInvalidRadius) {
		t.Fatalf("POIsWithinRadius() radius error = %v, want ErrInvalidRadius", err)
	}
}
