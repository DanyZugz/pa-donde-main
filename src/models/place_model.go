package models

import (
	"gorm.io/gorm"
)

type PlaceModel struct {
	gorm.Model
	Name             string
	ShortDescription string
	LongDescription  string
	Address          string

	Coordinates Coordinate  `gorm:"embedded"`
	PriceRange  PriceRange  `gorm:"embedded"`
	ContactInfo ContactInfo `gorm:"embedded"`
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type PriceRange struct {
	MinPrice float64
	MaxPrice float64
}

type ContactInfo struct {
	Email   string
	Phone   string
	Website string
}

// func GetPlacesByPriceRange(db *gorm.DB, minPrice float64, maxPrice float64) ([]PlaceModel, error) {
// 	var places []PlaceModel
// 	result := db.Where("min_price >= ? AND max_price <= ?", minPrice, maxPrice).Find(&places)
// 	return places, result.Error
// }
