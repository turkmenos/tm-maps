package tmmaps

import (
	"errors"
	"fmt"
	"math"
)

const earthRadiusKM = 6371.0088

// ErrNoSettlementCoordinates indicates that the bundled dataset contains no
// settlement with usable coordinates.
var ErrNoSettlementCoordinates = errors.New("no settlement coordinates available")

// NearestResult contains the nearest settlement and its great-circle distance
// from the requested coordinate in kilometres.
type NearestResult struct {
	Settlement
	DistanceKM float64 `json:"distance_km"`
}

// Nearest finds the closest known settlement to latitude and longitude using
// coordinates from the embedded dataset. DistanceKM is calculated with the
// Haversine formula.
func Nearest(latitude, longitude float64) (*NearestResult, error) {
	if !validNearestCoordinate(latitude, longitude) {
		return nil, fmt.Errorf("%w: latitude=%v longitude=%v", ErrInvalidCoordinate, latitude, longitude)
	}
	records, err := Regions()
	if err != nil {
		return nil, err
	}
	bySlug := make(map[string]*Region, len(records))
	for i := range records {
		bySlug[records[i].Slug] = &records[i]
	}

	var nearest *Region
	minimumDistance := math.Inf(1)
	for i := range records {
		record := &records[i]
		if !isSettlementType(record.Type) || record.Latitude == nil || record.Longitude == nil {
			continue
		}
		distance := haversineKM(latitude, longitude, *record.Latitude, *record.Longitude)
		if distance < minimumDistance {
			nearest = record
			minimumDistance = distance
		}
	}
	if nearest == nil {
		return nil, ErrNoSettlementCoordinates
	}
	return &NearestResult{
		Settlement: settlementFromRecord(nearest, bySlug),
		DistanceKM: minimumDistance,
	}, nil
}

func validNearestCoordinate(latitude, longitude float64) bool {
	return !math.IsNaN(latitude) && !math.IsNaN(longitude) &&
		!math.IsInf(latitude, 0) && !math.IsInf(longitude, 0) &&
		latitude >= -90 && latitude <= 90 && longitude >= -180 && longitude <= 180
}

func haversineKM(latitude1, longitude1, latitude2, longitude2 float64) float64 {
	lat1 := latitude1 * math.Pi / 180
	lat2 := latitude2 * math.Pi / 180
	deltaLat := (latitude2 - latitude1) * math.Pi / 180
	deltaLon := (longitude2 - longitude1) * math.Pi / 180

	sinLat := math.Sin(deltaLat / 2)
	sinLon := math.Sin(deltaLon / 2)
	a := sinLat*sinLat + math.Cos(lat1)*math.Cos(lat2)*sinLon*sinLon
	a = math.Min(1, math.Max(0, a))
	return earthRadiusKM * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
