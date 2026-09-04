package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"

	tmmaps "github.com/turkmenos/tm-maps"
)

const usage = `tm-maps exposes the bundled Turkmenistan map dataset.

Usage:
  tm-maps regions
  tm-maps find <slug>
  tm-maps children <parent-slug>
  tm-maps search [flags] <query>
  tm-maps nearest <latitude> <longitude>
  tm-maps within-radius <latitude> <longitude> <radius-km>
  tm-maps welayat <ahal|balkan|dasoguz|lebap|mary>
  tm-maps poi [category]
  tm-maps poi-search [flags] <query>
  tm-maps poi-nearby [flags] <latitude> <longitude> <radius-km>

Search flags:
  -limit int
  -region string
  -type string       may be repeated

POI categories:
  aeroway, cafes, hotels, pharmacies
`

type stringList []string

func (values *stringList) String() string {
	return fmt.Sprint([]string(*values))
}

func (values *stringList) Set(value string) error {
	*values = append(*values, value)
	return nil
}

// Run executes the tm-maps command-line interface.
func Run(stdout, stderr io.Writer, args []string) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}

	var value any
	var raw []byte
	var err error

	switch args[0] {
	case "regions":
		value, err = tmmaps.Regions()
	case "find":
		if len(args) != 2 {
			return usageError(stderr, "find requires exactly one slug")
		}
		value, err = tmmaps.FindRegion(args[1])
	case "children":
		if len(args) != 2 {
			return usageError(stderr, "children requires exactly one parent slug")
		}
		value, err = tmmaps.Children(args[1])
	case "search":
		value, err = runSearch(stderr, args[1:])
	case "nearest":
		value, err = runNearest(stderr, args[1:])
	case "within-radius":
		value, err = runWithinRadius(stderr, args[1:])
	case "welayat":
		if len(args) != 2 {
			return usageError(stderr, "welayat requires exactly one short slug")
		}
		raw, err = tmmaps.Welaýat(args[1])
	case "poi":
		value, err = runPOI(stderr, args[1:])
	case "poi-search":
		value, err = runPOISearch(stderr, args[1:])
	case "poi-nearby":
		value, err = runPOINearby(stderr, args[1:])
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return 0
	default:
		return usageError(stderr, "unknown command %q", args[0])
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	if raw != nil {
		_, _ = stdout.Write(raw)
		_, _ = fmt.Fprintln(stdout)
		return 0
	}
	if err := writeJSON(stdout, value); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	return 0
}

func runSearch(stderr io.Writer, args []string) (any, error) {
	flags := flag.NewFlagSet("search", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var types stringList
	options := tmmaps.SearchOptions{}
	flags.IntVar(&options.Limit, "limit", 0, "")
	flags.StringVar(&options.RegionSlug, "region", "", "")
	flags.Var(&types, "type", "")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if flags.NArg() != 1 {
		return nil, errors.New("search requires exactly one query")
	}
	options.Types = []string(types)
	return tmmaps.SearchWithOptions(flags.Arg(0), options)
}

func runNearest(stderr io.Writer, args []string) (any, error) {
	latitude, longitude, err := parseLatLon(args)
	if err != nil {
		return nil, err
	}
	return tmmaps.Nearest(latitude, longitude)
}

func runWithinRadius(stderr io.Writer, args []string) (any, error) {
	if len(args) != 3 {
		return nil, errors.New("within-radius requires latitude, longitude, and radius-km")
	}
	latitude, longitude, err := parseLatLon(args[:2])
	if err != nil {
		return nil, err
	}
	radius, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return nil, fmt.Errorf("parse radius-km: %w", err)
	}
	return tmmaps.WithinRadius(latitude, longitude, radius)
}

func runPOI(stderr io.Writer, args []string) (any, error) {
	switch len(args) {
	case 0:
		return tmmaps.AllPOIs()
	case 1:
		return tmmaps.POIs(tmmaps.POICategory(args[0]))
	default:
		return nil, errors.New("poi accepts at most one category")
	}
}

func runPOISearch(stderr io.Writer, args []string) (any, error) {
	flags := flag.NewFlagSet("poi-search", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var categories stringList
	flags.Var(&categories, "category", "")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if flags.NArg() != 1 {
		return nil, errors.New("poi-search requires exactly one query")
	}
	return tmmaps.SearchPOIs(flags.Arg(0), poiCategories(categories)...)
}

func runPOINearby(stderr io.Writer, args []string) (any, error) {
	flags := flag.NewFlagSet("poi-nearby", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var categories stringList
	flags.Var(&categories, "category", "")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if flags.NArg() != 3 {
		return nil, errors.New("poi-nearby requires latitude, longitude, and radius-km")
	}
	latitude, longitude, err := parseLatLon(flags.Args()[:2])
	if err != nil {
		return nil, err
	}
	radius, err := strconv.ParseFloat(flags.Arg(2), 64)
	if err != nil {
		return nil, fmt.Errorf("parse radius-km: %w", err)
	}
	return tmmaps.POIsWithinRadius(latitude, longitude, radius, poiCategories(categories)...)
}

func parseLatLon(args []string) (float64, float64, error) {
	if len(args) != 2 {
		return 0, 0, errors.New("requires latitude and longitude")
	}
	latitude, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse latitude: %w", err)
	}
	longitude, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse longitude: %w", err)
	}
	return latitude, longitude, nil
}

func poiCategories(values []string) []tmmaps.POICategory {
	categories := make([]tmmaps.POICategory, len(values))
	for i, value := range values {
		categories[i] = tmmaps.POICategory(value)
	}
	return categories
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func usageError(stderr io.Writer, format string, args ...any) int {
	fmt.Fprintf(stderr, "error: "+format+"\n\n", args...)
	fmt.Fprint(stderr, usage)
	return 2
}
