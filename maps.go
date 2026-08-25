// Package tmmaps provides Türkmenistan administrative boundary GeoJSON data.
package tmmaps

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

//go:embed data/regions.json data/geojson/welayatlar/*.geojson
var maps embed.FS
var (
	regionsOnce sync.Once
	regionsData []Region
	regionsErr  error
)

type Region struct {
	Slug               string   `json:"slug"`
	NameTM             string   `json:"name_tm"`
	NameEN             string   `json:"name_en"`
	NameRU             string   `json:"name_ru"`
	Type               string   `json:"type"`
	ParentSlug         string   `json:"parent_slug"`
	Latitude           *float64 `json:"latitude"`
	Longitude          *float64 `json:"longitude"`
	VerificationStatus string   `json:"verification_status"`
}

func Welaýat(name string) ([]byte, error) {
	switch name {
	case "ahal", "balkan", "dasoguz", "lebap", "mary":
	default:
		return nil, fmt.Errorf("unknown welaýat: %s", name)
	}
	return maps.ReadFile("data/geojson/welayatlar/" + name + ".geojson")
}
func Regions() ([]Region, error) {
	regionsOnce.Do(func() {
		data, err := maps.ReadFile("data/regions.json")
		if err != nil {
			regionsErr = fmt.Errorf("read regions: %w", err)
			return
		}
		if err := json.Unmarshal(data, &regionsData); err != nil {
			regionsErr = fmt.Errorf("decode regions: %w", err)
		}
	})
	if regionsErr != nil {
		return nil, regionsErr
	}
	return append([]Region(nil), regionsData...), nil
}
func FindRegion(slug string) (*Region, error) {
	regions, err := Regions()
	if err != nil {
		return nil, err
	}
	for i := range regions {
		if regions[i].Slug == slug {
			return &regions[i], nil
		}
	}
	return nil, fmt.Errorf("region not found: %s", slug)
}
func Children(parentSlug string) ([]Region, error) {
	regions, err := Regions()
	if err != nil {
		return nil, err
	}
	children := make([]Region, 0)
	for _, region := range regions {
		if region.ParentSlug == parentSlug {
			children = append(children, region)
		}
	}
	return children, nil
}
func Search(query string) ([]Region, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return []Region{}, nil
	}
	regions, err := Regions()
	if err != nil {
		return nil, err
	}
	results := make([]Region, 0)
	for _, region := range regions {
		if strings.Contains(strings.ToLower(region.NameTM), query) ||
			strings.Contains(strings.ToLower(region.NameEN), query) ||
			strings.Contains(strings.ToLower(region.NameRU), query) ||
			strings.Contains(strings.ToLower(region.Slug), query) {
			results = append(results, region)
		}
	}
	return results, nil
}
