package migrations

import (
	"strconv"
	"strings"
)

func parseUint(s string) uint {

	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)

	v, _ := strconv.ParseUint(s, 10, 64)

	return uint(v)
}

func parseInt(s string) int {

	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)

	v, _ := strconv.Atoi(s)

	return v
}

func parsePrice(s string) float64 {

	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)

	if s == "" {
		return 0
	}

	v, _ := strconv.ParseFloat(s, 64)

	// Karena format Indonesia 25.000 -> 25000
	return v
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func uintPtr(v uint) *uint {
	if v == 0 {
		return nil
	}
	return &v
}
