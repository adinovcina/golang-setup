package middleware

import (
	"net/http"

	"github.com/adinovcina/golang-setup/tools/paging"
)

type paginatorFetcher interface {
	SetPaginator(p paging.Paginator)
}

type paginatorCursorFetcher interface {
	SetPaginatorCursor(p paging.PaginatorCursor)
}

func Pagination(repo paginatorFetcher) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params := r.URL.Query()

			paginator := paging.NewPaginatorFromParams(params)

			repo.SetPaginator(paginator)

			// Delegate to next HTTP handler.
			next.ServeHTTP(w, r)
		})
	}
}

func PaginationCursor(repo paginatorCursorFetcher) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params := r.URL.Query()

			paginator := paging.NewPaginatorCursorFromParams(params)

			repo.SetPaginatorCursor(paginator)

			// Delegate to next HTTP handler.
			next.ServeHTTP(w, r)
		})
	}
}
