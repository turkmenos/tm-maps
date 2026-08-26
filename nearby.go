package tmmaps

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

// ErrInvalidRadius indicates a radius that is not finite and greater than zero.
var ErrInvalidRadius = errors.New("invalid radius")

// NearbySettlement contains a settlement and its distance from the requested
// coordinate in kilometres.
type NearbySettlement struct {
	Settlement
	DistanceKM float64 `json:"distance_km"`
}

// WithinRadius returns settlements no farther than radiusKM from latitude and
// longitude. Results are ordered from nearest to farthest.
func WithinRadius(latitude, longitude, radiusKM float64) ([]NearbySettlement, error) {
	if err := validateCoordinate(latitude, longitude); err != nil {
		return nil, err
	}
	if math.IsNaN(radiusKM) || math.IsInf(radiusKM, 0) || radiusKM <= 0 {
		return nil, fmt.Errorf("%w: %v km", ErrInvalidRadius, radiusKM)
	}

	records, err := Regions()
	if err != nil {
		return nil, err
	}
	bySlug := make(map[string]*Region, len(records))
	for i := range records {
		bySlug[records[i].Slug] = &records[i]
	}

	results := make([]NearbySettlement, 0)
	for i := range records {
		record := &records[i]
		if !isSettlementType(record.Type) || record.Latitude == nil || record.Longitude == nil {
			continue
		}
		distance := haversineKM(latitude, longitude, *record.Latitude, *record.Longitude)
		if distance <= radiusKM {
			results = append(results, NearbySettlement{
				Settlement: settlementFromRecord(record, bySlug),
				DistanceKM: distance,
			})
		}
	}
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].DistanceKM < results[j].DistanceKM
	})
	return results, nil
}
