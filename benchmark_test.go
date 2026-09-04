package tmmaps

import "testing"

func BenchmarkWelayat(b *testing.B) {
	for b.Loop() {
		if _, err := Welaýat("ahal"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRegions(b *testing.B) {
	for b.Loop() {
		if _, err := Regions(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFindRegion(b *testing.B) {
	for b.Loop() {
		if _, err := FindRegion("turkmenistan-mary"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkChildren(b *testing.B) {
	for b.Loop() {
		if _, err := Children("turkmenistan-mary"); err != nil {
			b.Fatal(err)
		}
	}
}

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

func BenchmarkAllPOIs(b *testing.B) {
	for b.Loop() {
		if _, err := AllPOIs(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPOIs(b *testing.B) {
	for b.Loop() {
		if _, err := POIs(POICategoryHotel); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearchPOIs(b *testing.B) {
	for b.Loop() {
		if _, err := SearchPOIs("dayanc", POICategoryHotel); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPOIsWithinRadius(b *testing.B) {
	for b.Loop() {
		if _, err := POIsWithinRadius(37.960077, 58.326063, 10, POICategoryCafe); err != nil {
			b.Fatal(err)
		}
	}
}
