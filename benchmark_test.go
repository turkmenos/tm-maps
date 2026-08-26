package tmmaps

import "testing"

func BenchmarkSearch(b *testing.B) {
	for b.Loop() {
		if _, err := Search("Mary"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearchWithOptions(b *testing.B) {
	options := SearchOptions{Limit: 10, Types: []string{"village"}, RegionSlug: "turkmenistan-mary"}
	for b.Loop() {
		if _, err := SearchWithOptions("a", options); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNearest(b *testing.B) {
	for b.Loop() {
		if _, err := Nearest(37.960077, 58.326063); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWithinRadius(b *testing.B) {
	for b.Loop() {
		if _, err := WithinRadius(37.960077, 58.326063, 100); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRegionAt(b *testing.B) {
	for b.Loop() {
		if _, err := RegionAt(37.960077, 58.326063); err != nil {
			b.Fatal(err)
		}
	}
}
