package api_v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pa-donde/src/models"
	"strconv"

	"gorm.io/gorm"
)

type PlaceHandler struct {
	DB *gorm.DB
}

func (p *PlaceHandler) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	// Fetch the places from the database
	var places []models.PlaceModel
	result := p.DB.Find(&places)

	if result.Error != nil {
		fmt.Println(result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Convert the places to JSON
	placesJSON, err := json.Marshal(places)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(placesJSON)
	
}

// GetPlaces returns a list of PlaceModel.
func (p *PlaceHandler) GetPlaces(w http.ResponseWriter, r *http.Request) {
	// Get the pagination values from the request context
	paginateDB, pagOk := r.Context().Value("paginateDB").(func(db *gorm.DB) *gorm.DB)

	if !pagOk {
		fmt.Println("Could not get paginateDB connection")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Fetch the places from the database
	var places []models.PlaceModel
	p.DB.Scopes(paginateDB).Find(&places)

	// json.NewEncoder(w).Encode(places)

	// Convert the places to JSON
	placesJSON, err := json.Marshal(places)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(placesJSON)
}

// GetPlace returns a single PlaceModel.
func (p *PlaceHandler) GetPlaceByID(w http.ResponseWriter, r *http.Request) {
	// Get the place ID from the URL parameters
	placeID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var place models.PlaceModel
	result := p.DB.First(&place, placeID)

	if result.Error != nil {
		fmt.Println(result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the place to JSON
	placeJSON, err := json.Marshal(place)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(placeJSON)
}

// CreatePlace creates a new PlaceModel.
func (p *PlaceHandler) CreatePlace(w http.ResponseWriter, r *http.Request) {
	// Get the place data from the request body
	var place models.PlaceModel
	err := json.NewDecoder(r.Body).Decode(&place)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create the place
	result := p.DB.Create(&place)
	if result.Error != nil {
		fmt.Println(result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the place to JSON
	placeJSON, err := json.Marshal(place)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(placeJSON)
}

// UpdatePlace updates an existing PlaceModel.
func (p *PlaceHandler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	// Get the place ID from the URL parameters
	_, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the place data from the request body
	var place models.PlaceModel

	err = json.NewDecoder(r.Body).Decode(&place)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the place
	result := p.DB.Model(&place).Updates(place)
	if result.Error != nil {
		fmt.Println(result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the place to JSON
	placeJSON, err := json.Marshal(place)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(placeJSON)
}

// DeletePlace deletes an existing PlaceModel.

func (p *PlaceHandler) DeletePlace(w http.ResponseWriter, r *http.Request) {
	// Get the place ID from the URL parameters
	placeID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Delete the place
	result := p.DB.Delete(&models.PlaceModel{}, placeID)
	if result.Error != nil {
		fmt.Println(result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
}
