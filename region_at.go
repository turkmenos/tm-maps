package tmmaps

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sync"
)

var (
	// ErrInvalidCoordinate indicates latitude or longitude outside its valid range.
	ErrInvalidCoordinate = errors.New("invalid coordinate")
	// ErrRegionNotFound indicates a coordinate outside the available welaýat boundaries.
	ErrRegionNotFound = errors.New("coordinate is outside available boundaries")
)

type boundaryGeometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

type boundaryFeature struct {
	Properties struct {
		Slug string `json:"slug"`
	} `json:"properties"`
	Geometry boundaryGeometry `json:"geometry"`
}

type boundaryCollection struct {
	Features []boundaryFeature `json:"features"`
}

var (
	boundariesOnce sync.Once
	boundariesData []boundaryFeature
	boundariesErr  error
)

// RegionAt returns the welaýat containing latitude and longitude.
// Boundary points are treated as contained. The lookup is entirely offline.
func RegionAt(latitude, longitude float64) (*Region, error) {
	if math.IsNaN(latitude) || math.IsNaN(longitude) ||
		math.IsInf(latitude, 0) || math.IsInf(longitude, 0) ||
		latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
		return nil, fmt.Errorf("%w: latitude=%v longitude=%v", ErrInvalidCoordinate, latitude, longitude)
	}

	boundaries, err := loadBoundaries()
	if err != nil {
		return nil, err
	}
	point := [2]float64{longitude, latitude}
	for _, feature := range boundaries {
		inside, err := geometryContains(feature.Geometry, point)
		if err != nil {
			return nil, fmt.Errorf("check boundary %s: %w", feature.Properties.Slug, err)
		}
		if inside {
			return FindRegion(feature.Properties.Slug)
		}
	}
	return nil, ErrRegionNotFound
}

func loadBoundaries() ([]boundaryFeature, error) {
	boundariesOnce.Do(func() {
		for _, name := range []string{"ahal", "balkan", "dasoguz", "lebap", "mary"} {
			data, err := Welaýat(name)
			if err != nil {
				boundariesErr = fmt.Errorf("read %s boundary: %w", name, err)
				return
			}
			var collection boundaryCollection
			if err := json.Unmarshal(data, &collection); err != nil {
				boundariesErr = fmt.Errorf("decode %s boundary: %w", name, err)
				return
			}
			boundariesData = append(boundariesData, collection.Features...)
		}
	})
	return boundariesData, boundariesErr
}

func geometryContains(geometry boundaryGeometry, point [2]float64) (bool, error) {
	switch geometry.Type {
	case "Polygon":
		var polygon [][][]float64
		if err := json.Unmarshal(geometry.Coordinates, &polygon); err != nil {
			return false, err
		}
		return polygonContains(polygon, point), nil
	case "MultiPolygon":
		var multiPolygon [][][][]float64
		if err := json.Unmarshal(geometry.Coordinates, &multiPolygon); err != nil {
			return false, err
		}
		for _, polygon := range multiPolygon {
			if polygonContains(polygon, point) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unsupported geometry type %q", geometry.Type)
	}
}

func polygonContains(polygon [][][]float64, point [2]float64) bool {
	if len(polygon) == 0 {
		return false
	}
	outer := ringLocation(polygon[0], point)
	if outer == 0 {
		return true
	}
	if outer < 0 {
		return false
	}
	for _, hole := range polygon[1:] {
		location := ringLocation(hole, point)
		if location == 0 {
			return true
		}
		if location > 0 {
			return false
		}
	}
	return true
}

// ringLocation returns -1 outside, 0 on the boundary, and 1 inside.
func ringLocation(ring [][]float64, point [2]float64) int {
	inside := false
	for i, j := 0, len(ring)-1; i < len(ring); j, i = i, i+1 {
		if len(ring[i]) < 2 || len(ring[j]) < 2 {
			continue
		}
		a := [2]float64{ring[j][0], ring[j][1]}
		b := [2]float64{ring[i][0], ring[i][1]}
		if pointOnSegment(point, a, b) {
			return 0
		}
		if (a[1] > point[1]) != (b[1] > point[1]) &&
			point[0] < (b[0]-a[0])*(point[1]-a[1])/(b[1]-a[1])+a[0] {
			inside = !inside
		}
	}
	if inside {
		return 1
	}
	return -1
}

func pointOnSegment(point, a, b [2]float64) bool {
	const epsilon = 1e-10
	lengthSquared := (b[0]-a[0])*(b[0]-a[0]) + (b[1]-a[1])*(b[1]-a[1])
	if lengthSquared <= epsilon*epsilon {
		return math.Abs(point[0]-a[0]) <= epsilon && math.Abs(point[1]-a[1]) <= epsilon
	}
	cross := (point[1]-a[1])*(b[0]-a[0]) - (point[0]-a[0])*(b[1]-a[1])
	if math.Abs(cross) > epsilon {
		return false
	}
	dot := (point[0]-a[0])*(b[0]-a[0]) + (point[1]-a[1])*(b[1]-a[1])
	if dot < -epsilon {
		return false
	}
	return dot <= lengthSquared+epsilon
}
