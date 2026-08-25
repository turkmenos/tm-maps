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
go get github.com/dayanchm/tm-maps
```

Usage:

```go id="2eq9pd"
package main

import (
	"fmt"

	tmmaps "github.com/dayanchm/tm-maps"
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