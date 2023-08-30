package utils

import (
	"pa-donde/src/interfaces"
	"pa-donde/src/middlewares"
	"pa-donde/src/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DEFAULT_PAGINATION_CONFIG = middlewares.PaginationConfig{
	DefaultPage:     1,
	DefaultPageSize: 10,
	MaxPageSize:     50,
}

const (
	TEST_DB_NAME = "test.db"
	PROD_DB_NAME = "pa-donde.db"
	SQLITE       = "sqlite"
	POSTGRES     = "postgres"
)

type Server struct {
	Router        *chi.Mux
	DB            *gorm.DB
	PlacesHandler interfaces.PlaceService
}

func CreateNewServer(db_name string, db_type string) *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	s.DB = configureDB(db_name, db_type)
	s.configureMiddleware()
	return s
}

func configureDB(db_name string, db_type string) *gorm.DB {
	// Create the database
	switch db_type {
	case SQLITE:
		db, err := gorm.Open(sqlite.Open(db_name), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		db.AutoMigrate(&models.PlaceModel{})
		return db
	case POSTGRES:
		panic("Postgres not implemented yet")
		// db, err := gorm.Open(postgres.Open(db_name), &gorm.Config{})
		// if err != nil {
		// 	panic("failed to connect database")
		// }
		// db.AutoMigrate(&models.PlaceModel{})
		// return db
	default:
		panic("Invalid database type")
	}

}

func createDB(db_name string) {
	// Create the database
	db, err := gorm.Open(sqlite.Open(db_name), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.PlaceModel{})
}

func (s *Server) configureMiddleware() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middlewares.SetDBMiddleware(s.DB))
}

func (s *Server) MountHandlers() {
	// Mount all handlers here
	s.Router.Route("/places", func(r chi.Router) {
		r.Get("/all", s.PlacesHandler.GetAllPlaces)
		r.With(middlewares.PaginationMiddleware(DEFAULT_PAGINATION_CONFIG)).Get("/{page}-{page_size}", s.PlacesHandler.GetPlaces)
		r.Get("/{placeID}", s.PlacesHandler.GetPlaceByID)
		r.Post("/", s.PlacesHandler.CreatePlace)
		r.Put("/{placeID}", s.PlacesHandler.UpdatePlace)
		r.Delete("/{placeID}", s.PlacesHandler.DeletePlace)

	})
}
