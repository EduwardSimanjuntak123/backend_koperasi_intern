package migrations

import (
	"strconv"
	"strings"
)

func parseUint(s string) uint {

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".0", "")

	f, _ := strconv.ParseFloat(s, 64)

	return uint(f)
}

func parseInt(s string) int {

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")

	i, _ := strconv.Atoi(s)

	return i
}

func parsePrice(s string) float64 {

	s = strings.TrimSpace(s)

	// 25.000 -> 25000
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")

	f, _ := strconv.ParseFloat(s, 64)

	return f
}
