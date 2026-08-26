# tm-maps

WGS84 (EPSG:4326) GeoJSON boundaries and geographic data for Turkmenistan, with a Go API for easy access.

**Version: v0.1**

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

## Go

Install:

```bash id="y9k7f3"
go get github.com/turkmenos/tm-maps
```

Usage:

```go id="2eq9pd"
package main

import (
	"fmt"

	tmmaps "github.com/turkmenos/tm-maps"
)

func main() {
	all, err := tmmaps.All()
	if err != nil {
		panic(err)
	}

	ahal, err := tmmaps.Welaýat("ahal")
	if err != nil {
		panic(err)
	}

	fmt.Println(len(all))
	fmt.Println(len(ahal))
}
```

The geographic data is embedded in the Go package, so no external files or network requests are required at runtime.

### Settlement search

`Search` finds settlements by their Turkmen, English, or Russian name. Matching
is case-insensitive, supports Turkmen Unicode characters, and works entirely
offline:

```go
results, err := tmmaps.Search("Mary")
if err != nil {
	panic(err)
}
for _, settlement := range results {
	fmt.Println(settlement.NameTM, settlement.Latitude, settlement.Longitude)
	if settlement.Region != nil {
		fmt.Println(settlement.Region.NameTM)
	}
}
```

Results are `Settlement` values and include coordinates when available, plus
the containing welaýat or independent-city information when assigned in the
source dataset. A query with no matches returns an empty slice.

### Nearest settlement

`Nearest` returns the closest settlement that has coordinates in the embedded
dataset. The result includes the Haversine distance in kilometres:

```go
place, err := tmmaps.Nearest(37.960077, 58.326063)
if err != nil {
	panic(err)
}
fmt.Println(place.NameTM)     // Aşgabat
fmt.Println(place.DistanceKM) // approximately 0
```

Invalid latitude or longitude values return `ErrInvalidCoordinate`.

### Coordinate lookup

`RegionAt` returns the welaýat containing a latitude and longitude:

```go
region, err := tmmaps.RegionAt(37.960077, 58.326063)
if err != nil {
	panic(err)
}
fmt.Println(region.NameTM) // Ahal
```

Invalid coordinates return `ErrInvalidCoordinate`. Coordinates outside the available
boundaries return `ErrRegionNotFound`; both can be checked with `errors.Is`.

The parameters use `(latitude, longitude)` order. GeoJSON coordinates are converted
internally from their standard `[longitude, latitude]` order.

`Contains` checks a coordinate against one specific welaýat. Points directly on
the boundary are considered contained:

```go
inside, err := tmmaps.Contains("ahal", 37.960077, 58.326063)
if err != nil {
	panic(err)
}
fmt.Println(inside) // true
```

An unsupported welaýat slug returns `ErrUnknownRegion`. Like `RegionAt`, this
operation uses only the boundary data embedded in the package.

## Data Coverage

The administrative boundaries are sourced from the geoBoundaries `TKM-ADM1-27578892` dataset and represent boundaries from **2007**.

Aşgabat is not represented as a separate ADM1 polygon in the source dataset and appears within Ahal. Arkadag, established as a city in 2023, is also not represented separately because the geometry dataset predates its creation.

For this reason, **v0.1 includes boundary geometries for the five welaýats:**

- Ahal
- Balkan
- Daşoguz
- Lebap
- Mary

Aşgabat and Arkadag may still appear in the geographic/settlement datasets where data is available, but they do not have separate ADM1 boundary geometries in this release.

## License

The source boundary dataset is released under the **Public Domain**.

See `DATA_LICENSE.md` for data licensing and attribution details.
