package tmmaps

import (
	"errors"
	"math"
	"testing"
)

func TestWithinRadius(t *testing.T) {
	results, err := WithinRadius(37.960077, 58.326063, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) < 2 {
		t.Fatalf("WithinRadius() returned %d results, want at least 2", len(results))
	}
	if results[0].Slug != "turkmenistan-asgabat" {
		t.Fatalf("WithinRadius() first slug = %q, want turkmenistan-asgabat", results[0].Slug)
	}
	for i := 1; i < len(results); i++ {
		if results[i].DistanceKM < results[i-1].DistanceKM {
			t.Fatalf("WithinRadius() results are not ordered at indexes %d and %d", i-1, i)
		}
		if results[i].DistanceKM > 100 {
			t.Fatalf("WithinRadius() distance = %v, want at most 100", results[i].DistanceKM)
		}
	}
}

func TestWithinRadiusEmpty(t *testing.T) {
	results, err := WithinRadius(0, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("WithinRadius() returned %d results, want none", len(results))
	}
}

func TestWithinRadiusInvalidInput(t *testing.T) {
	if _, err := WithinRadius(91, 0, 1); !errors.Is(err, ErrInvalidCoordinate) {
		t.Fatalf("WithinRadius() coordinate error = %v, want ErrInvalidCoordinate", err)
	}
	for _, radius := range []float64{0, -1, math.NaN(), math.Inf(1)} {
		if _, err := WithinRadius(0, 0, radius); !errors.Is(err, ErrInvalidRadius) {
			t.Errorf("WithinRadius() radius %v error = %v, want ErrInvalidRadius", radius, err)
		}
	}
}
