package cli

import (
	"io"
	"testing"
)

func BenchmarkRunSearch(b *testing.B) {
	for b.Loop() {
		if code := Run(io.Discard, io.Discard, []string{"search", "-limit", "10", "Mary"}); code != 0 {
			b.Fatalf("Run(search) exit code = %d, want 0", code)
		}
	}
}

func BenchmarkRunPOISearch(b *testing.B) {
	for b.Loop() {
		if code := Run(io.Discard, io.Discard, []string{"poi-search", "-category", "hotels", "dayanc"}); code != 0 {
			b.Fatalf("Run(poi-search) exit code = %d, want 0", code)
		}
	}
}

func BenchmarkRunPOINearby(b *testing.B) {
	for b.Loop() {
		if code := Run(io.Discard, io.Discard, []string{"poi-nearby", "-category", "cafes", "37.960077", "58.326063", "10"}); code != 0 {
			b.Fatalf("Run(poi-nearby) exit code = %d, want 0", code)
		}
	}
}
