package tmmaps

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/unicode/norm"
)

func normalizeSearchText(value string) string {
	return norm.NFC.String(cases.Fold().String(value))
}
