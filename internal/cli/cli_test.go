package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run(&stdout, &stderr, []string{"help"}); code != 0 {
		t.Fatalf("Run(help) exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "tm-maps regions") {
		t.Fatalf("Run(help) output missing usage: %q", stdout.String())
	}
}

func TestRunSearch(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run(&stdout, &stderr, []string{"search", "-limit", "1", "Mary"}); code != 0 {
		t.Fatalf("Run(search) exit code = %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"slug": "turkmenistan-mary"`) {
		t.Fatalf("Run(search) output missing Mary result: %s", stdout.String())
	}
}

func TestRunPOISearch(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run(&stdout, &stderr, []string{"poi-search", "-category", "hotels", "dayanc"}); code != 0 {
		t.Fatalf("Run(poi-search) exit code = %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"name": "Dayanc Hotel"`) {
		t.Fatalf("Run(poi-search) output missing hotel result: %s", stdout.String())
	}
}

func TestRunUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run(&stdout, &stderr, []string{"missing"}); code != 2 {
		t.Fatalf("Run(missing) exit code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("Run(missing) stderr missing error: %q", stderr.String())
	}
}
