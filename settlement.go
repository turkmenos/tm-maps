package tmmaps

import (
	"errors"
	"fmt"
	"strings"
)

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
	typeFilter, err := settlementTypeFilter(options.Types)
	if err != nil {
		return nil, err
	}

	query = normalizeSearchText(strings.TrimSpace(query))
	if query == "" {
		return []Settlement{}, nil
	}
	regions, err := Regions()
	if err != nil {
		return nil, err
	}
	bySlug := regionsBySlug(regions)

	results := make([]Settlement, 0)
	for i := range regions {
		record := &regions[i]
		if !settlementMatches(record, query, typeFilter) {
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

func settlementMatches(record *Region, query string, typeFilter map[string]struct{}) bool {
	if !isSettlementType(record.Type) {
		return false
	}
	if len(typeFilter) > 0 {
		if _, ok := typeFilter[record.Type]; !ok {
			return false
		}
	}
	return strings.Contains(normalizeSearchText(record.NameTM), query) ||
		strings.Contains(normalizeSearchText(record.NameEN), query) ||
		strings.Contains(normalizeSearchText(record.NameRU), query)
}

func settlementTypeFilter(types []string) (map[string]struct{}, error) {
	filter := make(map[string]struct{}, len(types))
	for _, settlementType := range types {
		if !isSettlementType(settlementType) {
			return nil, fmt.Errorf("%w: unknown settlement type %q", ErrInvalidSearchOptions, settlementType)
		}
		filter[settlementType] = struct{}{}
	}
	return filter, nil
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

func regionsBySlug(regions []Region) map[string]*Region {
	bySlug := make(map[string]*Region, len(regions))
	for i := range regions {
		bySlug[regions[i].Slug] = &regions[i]
	}
	return bySlug
}
