package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

// PaginationConfig holds the configuration for pagination.
type PaginationConfig struct {
	DefaultPage     int
	DefaultPageSize int
	MaxPageSize     int
}

// PaginationParams holds pagination parameters.
// type PaginationParams struct {
// 	Page     int
// 	PageSize int
// 	Offset   int
// }

// PaginationMiddleware is a middleware that adds pagination functionality to the endpoint.
func PaginationMiddleware(config PaginationConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract page and page size query parameters
			pageStr := r.URL.Query().Get("page")
			pageSizeStr := r.URL.Query().Get("page_size")

			// Convert page and page size parameters to integers
			page, err := strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				page = config.DefaultPage
			}

			pageSize, err := strconv.Atoi(pageSizeStr)
			if err != nil || pageSize < 1 {
				pageSize = config.DefaultPageSize
			}

			// Check if the page size exceeds the maximum allowed
			if pageSize > config.MaxPageSize {
				pageSize = config.MaxPageSize
			}

			// Calculate the offset and limit for fetching data
			offset := (page - 1) * pageSize

			db, ok := r.Context().Value("DB").(*gorm.DB)

			if !ok {
				fmt.Println("Could not get DB connection")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), "paginateDB", Paginate(db, offset, pageSize))
			r = r.WithContext(ctx)

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
		})
	}
}

func Paginate(db *gorm.DB, offset, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(pageSize)
	}
}
