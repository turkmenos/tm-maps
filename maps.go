// Package tmmaps provides Türkmenistan administrative boundary GeoJSON data.
package tmmaps

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/text/cases"
	"golang.org/x/text/unicode/norm"
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

// Settlement is a named populated place in the bundled geographic dataset.
// Latitude and Longitude are nil when coordinates are not available.
type Settlement struct {
	Slug      string            `json:"slug"`
	NameTM    string            `json:"name_tm"`
	NameEN    string            `json:"name_en"`
	NameRU    string            `json:"name_ru"`
	Type      string            `json:"type"`
	Latitude  *float64          `json:"latitude"`
	Longitude *float64          `json:"longitude"`
	Region    *SettlementRegion `json:"region,omitempty"`
}

// SettlementRegion identifies the top-level welaýat or independent city that
// contains a settlement. It is nil when the source record has no region assigned.
type SettlementRegion struct {
	Slug   string `json:"slug"`
	NameTM string `json:"name_tm"`
	NameEN string `json:"name_en"`
	NameRU string `json:"name_ru"`
	Type   string `json:"type"`
}

// SearchOptions controls settlement search filtering. A zero Limit means no
// limit. RegionSlug uses the full region slug, for example turkmenistan-mary.
type SearchOptions struct {
	Limit      int
	Types      []string
	RegionSlug string
}

// ErrInvalidSearchOptions indicates an invalid limit or settlement type.
var ErrInvalidSearchOptions = errors.New("invalid search options")

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

// Search finds settlements whose Turkmen, English, or Russian name contains
// query. Matching is case-insensitive and uses only the embedded dataset.
func Search(query string) ([]Settlement, error) {
	return SearchWithOptions(query, SearchOptions{})
}

// SearchWithOptions finds settlements by name and applies result, type, and
// region filters. Matching is case-insensitive and Unicode-normalized.
func SearchWithOptions(query string, options SearchOptions) ([]Settlement, error) {
	if options.Limit < 0 {
		return nil, fmt.Errorf("%w: limit must not be negative", ErrInvalidSearchOptions)
	}
	typeFilter := make(map[string]struct{}, len(options.Types))
	for _, settlementType := range options.Types {
		if !isSettlementType(settlementType) {
			return nil, fmt.Errorf("%w: unknown settlement type %q", ErrInvalidSearchOptions, settlementType)
		}
		typeFilter[settlementType] = struct{}{}
	}

	query = normalizeSearchText(strings.TrimSpace(query))
	if query == "" {
		return []Settlement{}, nil
	}
	regions, err := Regions()
	if err != nil {
		return nil, err
	}
	bySlug := make(map[string]*Region, len(regions))
	for i := range regions {
		bySlug[regions[i].Slug] = &regions[i]
	}

	results := make([]Settlement, 0)
	for i := range regions {
		record := &regions[i]
		if !isSettlementType(record.Type) {
			continue
		}
		if len(typeFilter) > 0 {
			if _, ok := typeFilter[record.Type]; !ok {
				continue
			}
		}
		if !strings.Contains(normalizeSearchText(record.NameTM), query) &&
			!strings.Contains(normalizeSearchText(record.NameEN), query) &&
			!strings.Contains(normalizeSearchText(record.NameRU), query) {
			continue
		}
		settlement := settlementFromRecord(record, bySlug)
		if options.RegionSlug != "" && (settlement.Region == nil || settlement.Region.Slug != options.RegionSlug) {
			continue
		}
		results = append(results, settlement)
		if options.Limit > 0 && len(results) == options.Limit {
			break
		}
	}
	return results, nil
}

func normalizeSearchText(value string) string {
	return norm.NFC.String(cases.Fold().String(value))
}

func settlementFromRecord(record *Region, bySlug map[string]*Region) Settlement {
	return Settlement{
		Slug:      record.Slug,
		NameTM:    record.NameTM,
		NameEN:    record.NameEN,
		NameRU:    record.NameRU,
		Type:      record.Type,
		Latitude:  record.Latitude,
		Longitude: record.Longitude,
		Region:    settlementRegion(record, bySlug),
	}
}

func isSettlementType(regionType string) bool {
	switch regionType {
	case "city", "town", "village", "independent_city":
		return true
	default:
		return false
	}
}

func settlementRegion(settlement *Region, bySlug map[string]*Region) *SettlementRegion {
	current := settlement
	for current != nil {
		if current.Type == "welayat" || current.Type == "independent_city" {
			return &SettlementRegion{
				Slug:   current.Slug,
				NameTM: current.NameTM,
				NameEN: current.NameEN,
				NameRU: current.NameRU,
				Type:   current.Type,
			}
		}
		current = bySlug[current.ParentSlug]
	}
	return nil
}
