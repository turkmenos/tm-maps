// Package tmmaps provides Türkmenistan administrative boundary GeoJSON data.
package tmmaps

import "embed"

//go:embed data/regions.json data/geojson/welayatlar/*.geojson data/poi/*.geojson
var maps embed.FS
