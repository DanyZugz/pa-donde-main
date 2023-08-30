package interfaces

import (
	"net/http"
)

type PlaceService interface {
	GetAllPlaces(w http.ResponseWriter, r *http.Request)
	GetPlaces(w http.ResponseWriter, r *http.Request)
	GetPlaceByID(w http.ResponseWriter, r *http.Request)
	CreatePlace(w http.ResponseWriter, r *http.Request)
	UpdatePlace(w http.ResponseWriter, r *http.Request)
	DeletePlace(w http.ResponseWriter, r *http.Request)
}
