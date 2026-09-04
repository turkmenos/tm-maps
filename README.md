# tm-maps

WGS84 (EPSG:4326) GeoJSON boundaries and geographic data for Turkmenistan, with a Go API for easy access.

## Live Demo

Try all public functions in the interactive web demo:

<https://tm-maps-demo-ten.vercel.app/>

Demo source: [`dayanchm/tm-maps-demo`](https://github.com/dayanchm/tm-maps-demo)

## GeoJSON

- `data/geojson/turkmenistan-welayatlar.geojson` — all welaýats in a single FeatureCollection
- `data/geojson/turkmenistan-etraplar.geojson` — available ADM2 district boundaries
- `data/geojson/yerlesim-noktalari.geojson` — 1,485 settlements with available coordinates
- `data/regions.json` — all 2,711 region and settlement records
- `data/regions/*.json` — records grouped by Ahal, Balkan, Daşoguz, Lebap, Mary, Aşgabat, and Arkadag
- `data/regions/unassigned.json` — 96 records that could not be assigned to a welaý
- `data/geojson/welayatlar/ahal.geojson`
- `data/geojson/welayatlar/balkan.geojson`
- `data/geojson/welayatlar/dasoguz.geojson`
- `data/geojson/welayatlar/lebap.geojson`
- `data/geojson/welayatlar/mary.geojson`

Each boundary feature contains:

- `slug`
- `name_tm`
- `name_en`
- `iso_3166_2`
- `admin_level`

GeoJSON coordinates follow the standard `[longitude, latitude]` order.

## Go API

Install the package:

```bash
go get github.com/turkmenos/tm-maps
```

Import it with:

```go
import tmmaps "github.com/turkmenos/tm-maps"
```

All data is embedded in the Go package. Every function works offline and does
not require external files, databases, or network services.

Geographic lookup functions accept coordinates in `(latitude, longitude)`
order. Raw GeoJSON follows the GeoJSON standard `[longitude, latitude]` order.

## CLI

The repository also includes a small command-line tool:

```bash
go run ./cmd/tm-maps help
go run ./cmd/tm-maps search -limit 5 Mary
go run ./cmd/tm-maps nearest 37.960077 58.326063
go run ./cmd/tm-maps poi-search -category hotels dayanc
go run ./cmd/tm-maps poi-nearby -category cafes 37.960077 58.326063 10
```

CLI output is formatted JSON, and it uses the same embedded offline dataset as
the Go API.

## Project Layout

```text
.
├── cmd/tm-maps/        # CLI entrypoint
├── internal/cli/       # CLI parsing and command handling
├── data/               # Embedded GeoJSON and region datasets
├── scripts/            # Data import and maintenance scripts
├── *.go                # Public tm-maps Go package
└── *_test.go           # Package and CLI tests
```

The root package exposes the reusable Go API. The command-line app stays under
`cmd/tm-maps`, while CLI-only implementation details live in `internal/cli`.
Bundled datasets stay under `data/` and are embedded by the package at build
time.

### `Welaýat`

```go
func Welaýat(name string) ([]byte, error)
```

Returns the raw GeoJSON `FeatureCollection` for one welaýat. Supported names
are `ahal`, `balkan`, `dasoguz`, `lebap`, and `mary`.

```go
data, err := tmmaps.Welaýat("ahal")
if err != nil {
	panic(err)
}
fmt.Println(string(data))
```

### `Regions`

```go
func Regions() ([]Region, error)
```

Returns all bundled country, welaýat, district, council, and settlement records
as structured `Region` values.

```go
regions, err := tmmaps.Regions()
if err != nil {
	panic(err)
}
for _, region := range regions {
	fmt.Println(region.Slug, region.NameTM, region.Type)
}
```

### `FindRegion`

```go
func FindRegion(slug string) (*Region, error)
```

Finds one record by its full slug.

```go
region, err := tmmaps.FindRegion("turkmenistan-mary")
if err != nil {
	panic(err)
}
fmt.Println(region.NameTM) // Mary
```

### `Children`

```go
func Children(parentSlug string) ([]Region, error)
```

Returns records whose direct parent matches `parentSlug`.

```go
children, err := tmmaps.Children("turkmenistan-dasoguz")
if err != nil {
	panic(err)
}
for _, child := range children {
	fmt.Println(child.NameTM, child.Type)
}
```

### `Search`

```go
func Search(query string) ([]Settlement, error)
```

Searches settlement names in Turkmen, English, and Russian. Matching is
case-insensitive and NFC-normalized. A query with no matches returns an empty
slice.

```go
results, err := tmmaps.Search("Mary")
if err != nil {
	panic(err)
}
for _, place := range results {
	fmt.Println(place.NameTM, place.Type)
	if place.Latitude != nil && place.Longitude != nil {
		fmt.Println(*place.Latitude, *place.Longitude)
	}
	if place.Region != nil {
		fmt.Println(place.Region.NameTM)
	}
}
```

Equivalent Unicode representations match the same settlement:

```go
results, _ := tmmaps.Search("Änew")
results, _ = tmmaps.Search("äNEW")
results, _ = tmmaps.Search("A\u0308new")
```

### `SearchWithOptions`

```go
func SearchWithOptions(query string, options SearchOptions) ([]Settlement, error)
```

Searches settlements with an optional result limit, settlement-type filter,
and region filter. A zero limit means unlimited results. Supported types are
`city`, `town`, `village`, and `independent_city`. `RegionSlug` uses a full slug.

```go
results, err := tmmaps.SearchWithOptions("a", tmmaps.SearchOptions{
	Limit:      10,
	Types:      []string{"city", "village"},
	RegionSlug: "turkmenistan-mary",
})
if err != nil {
	panic(err)
}
for _, place := range results {
	fmt.Println(place.NameTM)
}
```

A negative limit or unsupported type returns `ErrInvalidSearchOptions`.

### `RegionAt`

```go
func RegionAt(latitude, longitude float64) (*Region, error)
```

Returns the welaýat that covers a coordinate. Polygon boundary points count as
contained.

```go
region, err := tmmaps.RegionAt(37.960077, 58.326063)
if err != nil {
	panic(err)
}
fmt.Println(region.NameTM) // Ahal in the current boundary dataset
```

Invalid coordinates return `ErrInvalidCoordinate`. A coordinate outside the
available boundaries returns `ErrRegionNotFound`.

### `Contains`

```go
func Contains(slug string, latitude, longitude float64) (bool, error)
```

Reports whether a coordinate is covered by a specific welaýat. Use the short
slug: `ahal`, `balkan`, `dasoguz`, `lebap`, or `mary`. Boundary points return
`true`.

```go
inside, err := tmmaps.Contains("ahal", 37.960077, 58.326063)
if err != nil {
	panic(err)
}
fmt.Println(inside) // true
```

An unsupported slug returns `ErrUnknownRegion`. Invalid coordinates return
`ErrInvalidCoordinate`.

### `Nearest`

```go
func Nearest(latitude, longitude float64) (*NearestResult, error)
```

Returns the nearest settlement with known coordinates. `DistanceKM` contains
the Haversine distance in kilometres.

```go
place, err := tmmaps.Nearest(37.960077, 58.326063)
if err != nil {
	panic(err)
}
fmt.Println(place.NameTM)     // Aşgabat
fmt.Println(place.DistanceKM) // approximately 0
```

Invalid coordinates return `ErrInvalidCoordinate`. If no settlement coordinate
is available, the function returns `ErrNoSettlementCoordinates`.

### `WithinRadius`

```go
func WithinRadius(latitude, longitude, radiusKM float64) ([]NearbySettlement, error)
```

Returns settlements within `radiusKM`, ordered from nearest to farthest. Each
result includes `DistanceKM`. No matches produce an empty slice.

```go
places, err := tmmaps.WithinRadius(37.960077, 58.326063, 25)
if err != nil {
	panic(err)
}
for _, place := range places {
	fmt.Println(place.NameTM, place.DistanceKM)
}
```

Invalid coordinates return `ErrInvalidCoordinate`. A zero, negative, NaN, or
infinite radius returns `ErrInvalidRadius`.

### Result types

- `Region` represents any administrative or settlement record.
- `Settlement` represents a populated place and its optional coordinates and
  containing region.
- `SettlementRegion` describes the top-level welaýat or independent city for a
  settlement.
- `NearestResult` embeds `Settlement` and adds `DistanceKM`.
- `NearbySettlement` embeds `Settlement` and adds `DistanceKM`.
- `SearchOptions` provides `Limit`, `Types`, and `RegionSlug` filters.

Coordinate fields are pointers because some source records do not have known
coordinates. Check them for `nil` before dereferencing.

### Errors

Errors can be checked with `errors.Is`:

```go
place, err := tmmaps.Nearest(91, 0)
if errors.Is(err, tmmaps.ErrInvalidCoordinate) {
	fmt.Println("invalid latitude or longitude")
}
```

| Error | Meaning |
| --- | --- |
| `ErrInvalidCoordinate` | Latitude or longitude is outside its valid range, NaN, or infinite. |
| `ErrUnknownRegion` | `Contains` received an unsupported welaýat slug. |
| `ErrRegionNotFound` | A coordinate is outside the available boundary polygons. |
| `ErrNoSettlementCoordinates` | No settlement with usable coordinates is available. |
| `ErrInvalidRadius` | Radius is zero, negative, NaN, or infinite. |
| `ErrInvalidSearchOptions` | Search limit or settlement type is invalid. |

## Data Coverage

The administrative boundaries are sourced from the geoBoundaries `TKM-ADM1-27578892` dataset and represent boundaries from **2007**.

Aşgabat is not represented as a separate ADM1 polygon in the source dataset and appears within Ahal. Arkadag, established as a city in 2023, is also not represented separately because the geometry dataset predates its creation.

For this reason, the current dataset includes boundary geometries for five
welaýats:

- Ahal
- Balkan
- Daşoguz
- Lebap
- Mary

Aşgabat and Arkadag may still appear in the geographic/settlement datasets where data is available, but they do not have separate ADM1 boundary geometries in this release.

## License

The source boundary dataset is released under the **Public Domain**.

See `DATA_LICENSE.md` for data licensing and attribution details.
