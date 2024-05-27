package mysqlstore

import (
	"database/sql"

	"github.com/adinovcina/golang-setup/config"
	"github.com/adinovcina/golang-setup/store"
	"github.com/adinovcina/golang-setup/tools/paging"
)

const (
	ErrDuplicateEntry = 1062 // MySQL code for duplicate entry
)

// Compile-time check to assert implementation.
var _ store.Repository = (*Repository)(nil)

type Repository struct {
	db *sql.DB

	conf *config.Config

	paginator paging.Paginator

	paginatorCursor paging.PaginatorCursor
}

func New(db *sql.DB, conf *config.Config) *Repository {
	return &Repository{
		db:   db,
		conf: conf,
	}
}

func (r *Repository) SetPaginator(p paging.Paginator) {
	r.paginator = p
}

func (r *Repository) Paginator() paging.Paginator {
	if r.paginator == nil {
		return paging.NewPaginatorWithDefaults()
	}

	return r.paginator
}

func (r *Repository) SetPaginatorCursor(p paging.PaginatorCursor) {
	r.paginatorCursor = p
}

func (r *Repository) PaginatorCursor() paging.PaginatorCursor {
	if r.paginatorCursor == nil {
		return paging.NewPaginatorCursorWithDefaults()
	}

	return r.paginatorCursor
}
