package tmmaps

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
)

var (
	poiOnce sync.Once
	poiData []POI
	poiErr  error
)

// ErrUnknownPOICategory indicates that a requested point-of-interest category
// is not bundled with the package.
var ErrUnknownPOICategory = errors.New("unknown POI category")

// POICategory identifies a bundled point-of-interest category.
type POICategory string

const (
	POICategoryAeroway  POICategory = "aeroway"
	POICategoryCafe     POICategory = "cafes"
	POICategoryHotel    POICategory = "hotels"
	POICategoryPharmacy POICategory = "pharmacies"
)

// POI is a point of interest from the bundled OpenStreetMap extracts.
type POI struct {
	ID         string            `json:"id"`
	Category   POICategory       `json:"category"`
	Name       string            `json:"name,omitempty"`
	Latitude   float64           `json:"latitude"`
	Longitude  float64           `json:"longitude"`
	Properties map[string]string `json:"properties,omitempty"`
}

// NearbyPOI contains a point of interest and its distance from the requested
// coordinate in kilometres.
type NearbyPOI struct {
	POI
	DistanceKM float64 `json:"distance_km"`
}

type poiCollection struct {
	Features []poiFeature `json:"features"`
}

type poiFeature struct {
	ID         string            `json:"id"`
	Properties map[string]string `json:"properties"`
	Geometry   struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
	} `json:"geometry"`
}

// POIs returns all bundled points of interest for category. Supported
// categories are aeroway, cafes, hotels, and pharmacies.
func POIs(category POICategory) ([]POI, error) {
	if !isPOICategory(category) {
		return nil, fmt.Errorf("%w: %s", ErrUnknownPOICategory, category)
	}
	points, err := loadPOIs()
	if err != nil {
		return nil, err
	}
	results := make([]POI, 0)
	for _, point := range points {
		if point.Category == category {
			results = append(results, copyPOI(point))
		}
	}
	return results, nil
}

// AllPOIs returns every bundled point of interest.
func AllPOIs() ([]POI, error) {
	points, err := loadPOIs()
	if err != nil {
		return nil, err
	}
	results := make([]POI, len(points))
	for i, point := range points {
		results[i] = copyPOI(point)
	}
	return results, nil
}

// SearchPOIs finds bundled points of interest whose names contain query.
// Matching is case-insensitive and Unicode-normalized.
func SearchPOIs(query string, categories ...POICategory) ([]POI, error) {
	categoryFilter, err := poiCategoryFilter(categories)
	if err != nil {
		return nil, err
	}
	query = normalizeSearchText(strings.TrimSpace(query))
	if query == "" {
		return []POI{}, nil
	}
	points, err := loadPOIs()
	if err != nil {
		return nil, err
	}
	results := make([]POI, 0)
	for _, point := range points {
		if len(categoryFilter) > 0 {
			if _, ok := categoryFilter[point.Category]; !ok {
				continue
			}
		}
		if strings.Contains(normalizeSearchText(point.Name), query) {
			results = append(results, copyPOI(point))
		}
	}
	return results, nil
}

// POIsWithinRadius returns points of interest within radiusKM, ordered from
// nearest to farthest. When categories is empty, all POI categories are used.
func POIsWithinRadius(latitude, longitude, radiusKM float64, categories ...POICategory) ([]NearbyPOI, error) {
	if err := validateCoordinate(latitude, longitude); err != nil {
		return nil, err
	}
	if math.IsNaN(radiusKM) || math.IsInf(radiusKM, 0) || radiusKM <= 0 {
		return nil, fmt.Errorf("%w: %v km", ErrInvalidRadius, radiusKM)
	}
	categoryFilter, err := poiCategoryFilter(categories)
	if err != nil {
		return nil, err
	}
	points, err := loadPOIs()
	if err != nil {
		return nil, err
	}
	results := make([]NearbyPOI, 0)
	for _, point := range points {
		if len(categoryFilter) > 0 {
			if _, ok := categoryFilter[point.Category]; !ok {
				continue
			}
		}
		distance := haversineKM(latitude, longitude, point.Latitude, point.Longitude)
		if distance <= radiusKM {
			results = append(results, NearbyPOI{
				POI:        copyPOI(point),
				DistanceKM: distance,
			})
		}
	}
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].DistanceKM < results[j].DistanceKM
	})
	return results, nil
}

func loadPOIs() ([]POI, error) {
	poiOnce.Do(func() {
		for _, category := range []POICategory{POICategoryAeroway, POICategoryCafe, POICategoryHotel, POICategoryPharmacy} {
			data, err := maps.ReadFile("data/poi/" + string(category) + ".geojson")
			if err != nil {
				poiErr = fmt.Errorf("read %s POIs: %w", category, err)
				return
			}
			var collection poiCollection
			if err := json.Unmarshal(data, &collection); err != nil {
				poiErr = fmt.Errorf("decode %s POIs: %w", category, err)
				return
			}
			for i, feature := range collection.Features {
				latitude, longitude, err := poiCoordinate(feature.Geometry.Type, feature.Geometry.Coordinates)
				if err != nil {
					poiErr = fmt.Errorf("decode %s POIs: feature %d: %w", category, i, err)
					return
				}
				if err := validateCoordinate(latitude, longitude); err != nil {
					poiErr = fmt.Errorf("decode %s POIs: feature %d: %w", category, i, err)
					return
				}
				poiData = append(poiData, POI{
					ID:         firstNonEmpty(feature.ID, feature.Properties["@id"]),
					Category:   category,
					Name:       feature.Properties["name"],
					Latitude:   latitude,
					Longitude:  longitude,
					Properties: copyStringMap(feature.Properties),
				})
			}
		}
	})
	return poiData, poiErr
}

func poiCategoryFilter(categories []POICategory) (map[POICategory]struct{}, error) {
	filter := make(map[POICategory]struct{}, len(categories))
	for _, category := range categories {
		if !isPOICategory(category) {
			return nil, fmt.Errorf("%w: %s", ErrUnknownPOICategory, category)
		}
		filter[category] = struct{}{}
	}
	return filter, nil
}

func isPOICategory(category POICategory) bool {
	switch category {
	case POICategoryAeroway, POICategoryCafe, POICategoryHotel, POICategoryPharmacy:
		return true
	default:
		return false
	}
}

func poiCoordinate(geometryType string, coordinates json.RawMessage) (float64, float64, error) {
	bounds := poiBounds{
		minLat: math.Inf(1),
		maxLat: math.Inf(-1),
		minLon: math.Inf(1),
		maxLon: math.Inf(-1),
	}
	switch geometryType {
	case "Point":
		var point []float64
		if err := json.Unmarshal(coordinates, &point); err != nil {
			return 0, 0, err
		}
		if len(point) < 2 {
			return 0, 0, errors.New("point has fewer than two coordinates")
		}
		return point[1], point[0], nil
	case "Polygon":
		var polygon [][][]float64
		if err := json.Unmarshal(coordinates, &polygon); err != nil {
			return 0, 0, err
		}
		bounds.addPolygon(polygon)
	case "MultiPolygon":
		var multiPolygon [][][][]float64
		if err := json.Unmarshal(coordinates, &multiPolygon); err != nil {
			return 0, 0, err
		}
		for _, polygon := range multiPolygon {
			bounds.addPolygon(polygon)
		}
	default:
		return 0, 0, fmt.Errorf("unsupported geometry type %q", geometryType)
	}
	if !bounds.valid() {
		return 0, 0, errors.New("geometry has no usable coordinates")
	}
	return (bounds.minLat + bounds.maxLat) / 2, (bounds.minLon + bounds.maxLon) / 2, nil
}

type poiBounds struct {
	minLat float64
	maxLat float64
	minLon float64
	maxLon float64
}

func (bounds *poiBounds) addPolygon(polygon [][][]float64) {
	for _, ring := range polygon {
		for _, coordinate := range ring {
			if len(coordinate) < 2 {
				continue
			}
			longitude := coordinate[0]
			latitude := coordinate[1]
			bounds.minLat = math.Min(bounds.minLat, latitude)
			bounds.maxLat = math.Max(bounds.maxLat, latitude)
			bounds.minLon = math.Min(bounds.minLon, longitude)
			bounds.maxLon = math.Max(bounds.maxLon, longitude)
		}
	}
}

func (bounds poiBounds) valid() bool {
	return !math.IsInf(bounds.minLat, 0) && !math.IsInf(bounds.maxLat, 0) &&
		!math.IsInf(bounds.minLon, 0) && !math.IsInf(bounds.maxLon, 0)
}

func copyPOI(point POI) POI {
	point.Properties = copyStringMap(point.Properties)
	return point
}
