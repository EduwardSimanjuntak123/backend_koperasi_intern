package migrations

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"backend_koperasi/internal/models"
	"backend_koperasi/internal/utils"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func SeedProductsFromExcel(db *gorm.DB, filePath string) error {

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	sheet := f.GetSheetName(0)

	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}
	previewLimit := min(10, len(rows))
	for i := 0; i < previewLimit; i++ {
		fmt.Println(rows[i])
	}

	if len(rows) == 0 {
		return errors.New("excel kosong")
	}

	// ============================
	// Cari Header
	// ============================

	headerIndex := -1

	for i, row := range rows {

		if len(row) < 5 {
			continue
		}

		hasNo := false
		hasBarcode := false
		hasName := false

		for _, cell := range row {
			v := strings.TrimSpace(cell)

			if strings.EqualFold(v, "No") {
				hasNo = true
			}

			if strings.EqualFold(v, "Barcode") {
				hasBarcode = true
			}

			if strings.EqualFold(v, "Name") {
				hasName = true
			}
		}

		if hasNo && hasBarcode && hasName {
			headerIndex = i
			break
		}
	}

	if headerIndex == -1 {
		return errors.New("header excel tidak ditemukan")
	}

	headerRow := rows[headerIndex]
	noCol := findColumnIndex(headerRow, "No")
	if noCol == -1 {
		return errors.New("kolom 'No' tidak ditemukan pada header")
	}

	fmt.Println("Header ditemukan pada row", headerIndex+1)
	fmt.Println("Kolom 'No' berada di index", noCol)

	success := 0
	skipped := 0
	failed := 0

	// ============================
	// Import
	// ============================

	for i := headerIndex + 1; i < len(rows); i++ {

		row := rows[i]

		if len(row) < noCol+17 {
			continue
		}

		col := func(offset int) string {
			return getCell(row, noCol+offset)
		}

		name := strings.TrimSpace(col(4))
		barcode := strings.TrimSpace(col(2))

		// Skip repeated header rows or non-data rows.
		if strings.EqualFold(name, "Name") || strings.EqualFold(barcode, "Barcode") {
			continue
		}

		if name == "" {
			continue
		}

		if barcode == "" {
			continue
		}

		var exist models.Product

		err := db.Where("barcode = ?", barcode).First(&exist).Error

		if err == nil {
			fmt.Println(barcode, "sudah ada")
			skipped++
			continue
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		product := models.Product{}

		// =========================
		// Basic
		// =========================

		product.Name = name
		product.Slug = utils.GenerateSlug(name)

		product.Barcode = barcode
		product.SKU = col(3)

		image := strings.TrimSpace(col(1))
		if image != "" {
			product.Image = &image
		}

		// =========================
		// Relation
		// =========================

		category := parseUint(col(5))
		if category != 0 {
			product.CategoryID = &category
		}

		brand := parseUint(col(8))
		if brand != 0 {
			product.BrandID = &brand
		}

		unit := parseUint(col(9))
		if unit != 0 {
			product.UnitID = &unit
		}

		product.StoreID = parseUint(col(7))

		// =========================
		// Inventory
		// =========================

		mv := models.FastMoving

		if strings.EqualFold(col(6), "Slow Moving") {
			mv = models.SlowMoving
		}

		product.InventoryMovement = &mv

		// =========================
		// Price
		// =========================

		product.Price = parsePrice(col(10))
		product.PurchasePrice = parsePrice(col(11))

		promo := parsePrice(col(12))
		if promo > 0 {
			product.PromotionPrice = &promo
		}

		// =========================
		// Stock
		// =========================

		product.Stock = parseInt(col(16))

		min := parseInt(col(13))
		if min > 0 {
			product.MinStock = &min
		}

		max := parseInt(col(14))
		if max > 0 {
			product.MaxStock = &max
		}

		// =========================
		// Expired
		// =========================

		expired := strings.TrimSpace(col(15))

		if expired != "" {
			if t, ok := parseExcelDate(expired); ok {
				product.ExpiredDate = &t
			}
		}

		// =========================
		// Insert
		// =========================

		if err := db.Create(&product).Error; err != nil {

			failed++

			fmt.Printf(
				"Baris %d gagal : %v\n",
				i+1,
				err,
			)

			continue
		}

		success++
	}

	fmt.Println("==========================")
	fmt.Println("IMPORT SELESAI")
	fmt.Println("==========================")
	fmt.Println("Berhasil :", success)
	fmt.Println("Sudah ada :", skipped)
	fmt.Println("Gagal :", failed)

	return nil
}

func getCell(row []string, index int) string {

	if index >= len(row) {
		return ""
	}

	return strings.TrimSpace(row[index])
}

func findColumnIndex(row []string, headerName string) int {
	for i, cell := range row {
		if strings.EqualFold(strings.TrimSpace(cell), headerName) {
			return i
		}
	}

	return -1
}

func cleanNumber(s string) string {

	s = strings.TrimSpace(s)

	s = strings.ReplaceAll(s, "\u00A0", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")

	return s
}

func parseExcelDate(raw string) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false
	}

	layouts := []string{
		"02/01/2006",
		"02-01-2006",
		"02/01/06",
		"02-01-06",
		"2006-01-02",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, raw)
		if err == nil {
			return t, true
		}
	}

	return time.Time{}, false
}
