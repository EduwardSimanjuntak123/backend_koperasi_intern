package utils

import (
	"regexp"
	"strings"
)

func GenerateSlug(text string) string {
	// lowercase
	slug := strings.ToLower(text)

	// hapus spasi di depan dan belakang
	slug = strings.TrimSpace(slug)

	// ganti semua spasi menjadi "-"
	slug = regexp.MustCompile(`\s+`).ReplaceAllString(slug, "-")

	// hapus karakter selain huruf, angka, dan "-"
	slug = regexp.MustCompile(`[^a-z0-9\-]`).ReplaceAllString(slug, "")

	// gabungkan "-" yang berulang
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	// hapus "-" di awal/akhir
	slug = strings.Trim(slug, "-")

	return slug
}
