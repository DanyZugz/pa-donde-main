// main_test.go
package api_v1

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"pa-donde/src/models"
	"pa-donde/src/utils"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *utils.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func TestPlacesTableDriven(t *testing.T) {
	s := utils.CreateNewServer(utils.TEST_DB_NAME, utils.SQLITE)
	s.PlacesHandler = &PlaceHandler{DB: s.DB}
	s.MountHandlers()

	// Test setup and teardown
	setupTestData(t, s.DB)
	defer cleanupTestData(t, s.DB)

	// Table-driven tests
	testCases := []struct {
		name           string
		method         string
		url            string
		payload        interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "CreatePlace",
			method:         "POST",
			url:            "/places",
			payload:        getRandomTestPlace(),
			expectedStatus: http.StatusCreated,
			expectedBody:   "",
		},
		{
			name:           "GetAllPlaces",
			method:         "GET",
			url:            "/places/all",
			payload:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "GetPlaceByID",
			method:         "GET",
			url:            "/places?placeID=" + strconv.Itoa(int(getTestPlaceID(t, s.DB))),
			payload:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "UpdatePlace",
			method:         "PUT",
			url:            "/places?placeID=" + strconv.Itoa(int(getTestPlaceID(t, s.DB))),
			payload:        getRandomTestPlace(),
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "DeletePlace",
			method:         "DELETE",
			url:            "/places?placeID=" + strconv.Itoa(int(getTestPlaceID(t, s.DB))),
			payload:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest(tc.method, tc.url, bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req, s)

			assert.Equal(t, tc.expectedStatus, response.Code)
			assert.Contains(t, response.Body.String(), tc.expectedBody)

		})
	}
}

// func setupServer(t *testing.T) *Server {
// 	r := chi.NewRouter()
// 	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
// 	if err != nil {
// 		t.Fatalf("failed to connect to the test database: %v", err)
// 	}

// 	// Auto-migrate the test model to create the "places" table.
// 	db.AutoMigrate(&models.PlaceModel{})
// 	setupTestData(t, db)
// 	return s
// }

// func setupServer(t *testing.T) *gorm.DB {

// 	return db
// }

func setupTestData(t *testing.T, db *gorm.DB) {
	// Seed the test database with some initial data for testing
	testPlace := getTestPlace()
	if err := db.Create(&testPlace).Error; err != nil {
		t.Fatalf("failed to seed test data: %v", err)
	}
}

func cleanupTestData(t *testing.T, db *gorm.DB) {
	// Cleanup the test database after all test cases finish
	if err := db.Exec("DELETE FROM place_models").Error; err != nil {
		t.Fatalf("failed to clean up test data: %v", err)
	}
}

func getTestPlace() models.PlaceModel {
	return models.PlaceModel{
		Name:             "Test Place",
		ShortDescription: "Short description",
		LongDescription:  "Long description",
		Address:          "Test Address",
		Coordinates: models.Coordinate{
			Latitude:  12.345,
			Longitude: 67.890,
		},
		PriceRange: models.PriceRange{
			MinPrice: 5.0,
			MaxPrice: 20.0,
		},
		ContactInfo: models.ContactInfo{
			Email:   "test@example.com",
			Phone:   "123456789",
			Website: "http://example.com",
		},
	}
}

// getRandomTestPlace generates a random test PlaceModel for each test run.
func getRandomTestPlace() models.PlaceModel {
	// Initialize the random seed based on the current time.
	rand.Seed(time.Now().UnixNano())

	return models.PlaceModel{
		Name:             randomString(10),
		ShortDescription: randomString(20),
		LongDescription:  randomString(50),
		Address:          randomString(30),
		Coordinates: models.Coordinate{
			Latitude:  rand.Float64()*180.0 - 90.0,
			Longitude: rand.Float64()*360.0 - 180.0,
		},
		PriceRange: models.PriceRange{
			MinPrice: rand.Float64() * 50.0,
			MaxPrice: rand.Float64() * 100.0,
		},
		ContactInfo: models.ContactInfo{
			Email:   randomString(10) + "@example.com",
			Phone:   randomNumericString(10),
			Website: "http://" + randomString(10) + ".com",
		},
	}
}

func getTestPlaceID(t *testing.T, db *gorm.DB) uint {
	var place models.PlaceModel

	result := db.First(&place)

	if result.Error != nil {
		t.Fatalf("Failed to get test place ID: %v", result.Error)
	}

	return place.ID
}

// randomString generates a random alphanumeric string of given length.
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// randomNumericString generates a random numeric string of given length.
func randomNumericString(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
