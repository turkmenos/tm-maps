package tmmaps

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestWelaýat(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"ahal", false},
		{"ashgabat", true},
		{"balkan", false},
		{"dasoguz", false},
		{"lebap", false},
		{"mary", false},
		{"invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Welaýat(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Welaýat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				var geojson map[string]interface{}
				if err := json.Unmarshal(data, &geojson); err != nil {
					t.Errorf("Welaýat() returned invalid GeoJSON: %v", err)
				}
			}
		})
	}
}

func TestRegions(t *testing.T) {
	regions, err := Regions()
	if err != nil {
		t.Fatal(err)
	}
	if len(regions) != 2711 {
		t.Fatalf("Regions() returned %d records, want 2711", len(regions))
	}
}

func TestFindRegion(t *testing.T) {
	region, err := FindRegion("turkmenistan-mary")
	if err != nil {
		t.Fatal(err)
	}
	if region.NameTM != "Mary" || region.Type != "welayat" {
		t.Fatalf("FindRegion() returned %+v", region)
	}

	if _, err := FindRegion("does-not-exist"); err == nil {
		t.Fatal("FindRegion() expected an error for an unknown slug")
	}
}

func TestChildren(t *testing.T) {
	children, err := Children("turkmenistan-dasoguz")
	if err != nil {
		t.Fatal(err)
	}
	if len(children) == 0 {
		t.Fatal("Children() returned no Daşoguz children")
	}
	for _, child := range children {
		if child.ParentSlug != "turkmenistan-dasoguz" {
			t.Fatalf("unexpected child parent: %s", child.ParentSlug)
		}
	}
}

func TestSearch(t *testing.T) {
	results, err := Search("Mary")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("Search() returned no results for Mary")
	}
	for _, result := range results {
		if !isSettlementType(result.Type) {
			t.Fatalf("Search() returned non-settlement type %q", result.Type)
		}
		if result.Region == nil || result.Region.Slug != "turkmenistan-mary" {
			t.Fatalf("Search() result has unexpected region: %+v", result.Region)
		}
	}

	unicodeResults, err := Search("äNEW")
	if err != nil {
		t.Fatal(err)
	}
	if len(unicodeResults) == 0 || unicodeResults[0].NameTM != "Änew" {
		t.Fatalf("Search() did not match Turkmen Unicode case-insensitively: %+v", unicodeResults)
	}
	decomposedResults, err := Search("A\u0308new")
	if err != nil {
		t.Fatal(err)
	}
	if len(decomposedResults) == 0 || decomposedResults[0].NameTM != "Änew" {
		t.Fatalf("Search() did not normalize Unicode: %+v", decomposedResults)
	}

	coordinateResults, err := Search("aşgabat")
	if err != nil {
		t.Fatal(err)
	}
	if len(coordinateResults) == 0 || coordinateResults[0].Latitude == nil || coordinateResults[0].Longitude == nil {
		t.Fatalf("Search() did not include available coordinates: %+v", coordinateResults)
	}

	empty, err := Search("no-such-settlement-xyz")
	if err != nil {
		t.Fatal(err)
	}
	if len(empty) != 0 {
		t.Fatalf("Search() returned %d results for an unknown query", len(empty))
	}
}

func TestSearchWithOptions(t *testing.T) {
	results, err := SearchWithOptions("a", SearchOptions{
		Limit:      3,
		Types:      []string{"village"},
		RegionSlug: "turkmenistan-mary",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("SearchWithOptions() returned %d results, want 3", len(results))
	}
	for _, result := range results {
		if result.Type != "village" {
			t.Fatalf("SearchWithOptions() type = %q, want village", result.Type)
		}
		if result.Region == nil || result.Region.Slug != "turkmenistan-mary" {
			t.Fatalf("SearchWithOptions() region = %+v, want turkmenistan-mary", result.Region)
		}
	}
}

func TestSearchWithOptionsInvalid(t *testing.T) {
	for _, options := range []SearchOptions{{Limit: -1}, {Types: []string{"etrap"}}} {
		if _, err := SearchWithOptions("Mary", options); !errors.Is(err, ErrInvalidSearchOptions) {
			t.Errorf("SearchWithOptions() error = %v, want ErrInvalidSearchOptions", err)
		}
	}
}
