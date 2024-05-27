package store

import "github.com/adinovcina/golang-setup/tools/paging"

// Data Access Layer (DAL) methods.
type Repository interface {
	// Pagination
	Paginator() paging.Paginator
	SetPaginator(p paging.Paginator)
	// Pagination cursor
	PaginatorCursor() paging.PaginatorCursor
	SetPaginatorCursor(p paging.PaginatorCursor)

	AccountRepository
	UserRepository
	TokenRepository
}

type InMemRepository interface {
	AccountInMemRepository
}
