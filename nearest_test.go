package tmmaps

import (
	"errors"
	"math"
	"testing"
)

func TestNearestKnownLocation(t *testing.T) {
	result, err := Nearest(37.960077, 58.326063)
	if err != nil {
		t.Fatal(err)
	}
	if result.Slug != "turkmenistan-asgabat" {
		t.Fatalf("Nearest() slug = %q, want turkmenistan-asgabat", result.Slug)
	}
	if result.DistanceKM > 0.001 {
		t.Fatalf("Nearest() distance = %v km, want approximately zero", result.DistanceKM)
	}
	if result.Latitude == nil || result.Longitude == nil {
		t.Fatal("Nearest() result omitted settlement coordinates")
	}
	if result.Region == nil || result.Region.Slug != "turkmenistan-asgabat" {
		t.Fatalf("Nearest() result has unexpected region: %+v", result.Region)
	}
}

func TestNearestInvalidCoordinate(t *testing.T) {
	_, err := Nearest(91, 0)
	if !errors.Is(err, ErrInvalidCoordinate) {
		t.Fatalf("Nearest() error = %v, want ErrInvalidCoordinate", err)
	}
}

func TestHaversineKM(t *testing.T) {
	got := haversineKM(0, 0, 0, 1)
	if math.Abs(got-111.195) > 0.001 {
		t.Fatalf("haversineKM() = %v, want approximately 111.195", got)
	}
}
