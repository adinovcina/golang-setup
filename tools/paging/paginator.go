package paging

import (
	"strconv"
	"strings"
)

const (
	// paginatorLimitDefault - Default amount of results to be returned.
	paginatorLimitDefault = 20
	// paginatorPageDefault - Default page number.
	paginatorPageDefault = 1
	// paginatorDefaultOrderDir - Default order direction.
	paginatorDefaultOrderDir = "desc"
	// paginatorDefaultSort - Default order by which is an empty string.
	paginatorDefaultSort = ""
)

type Paginator interface {
	Save(count, responseLength int)
	GetLimit() int
	GetOffset() int
	Order(defaultSortKey, defaultDir string) (key, direction string)
}

// PaginatorParams is a parameters provider interface to get the pagination params from.
type PaginatorParams interface {
	Get(key string) string
}

type Pagination struct {
	// Sort - Sort field in a format - sort=field1
	Sort string `json:"sort"`
	// Limit - A limit on the number of objects to be returned
	Limit int `json:"limit"`
	// Offset - Page * Limit (ex: 2 * 10, Offset == 20)
	Offset int `json:"offset"`
	// Page - Current page you're on
	Page int `json:"page"`
	// TotalEntriesSize - Total potential records matching the query
	TotalEntriesSize int `json:"totalEntriesSize"`
	// CurrentEntriesSize - Total records returned, will be <= Limit
	CurrentEntriesSize int `json:"currentEntriesSize"`
	// TotalPages - Number of total pages
	TotalPages int `json:"totalPages"`
}

func NewPaginatorWithDefaults() *Pagination {
	return NewPaginator(paginatorPageDefault, paginatorLimitDefault, paginatorDefaultSort)
}

func NewPaginator(page, limit int, sort string) *Pagination {
	if page < 1 {
		page = paginatorPageDefault
	}

	if limit < 1 {
		limit = paginatorLimitDefault
	}

	paginator := &Pagination{
		Limit:              limit,
		Page:               page,
		Offset:             0,
		TotalEntriesSize:   0,
		CurrentEntriesSize: 0,
		TotalPages:         0,
		Sort:               "",
	}
	paginator.Offset = (page - 1) * limit

	if sort != "" {
		paginator.Sort = sort
	}

	return paginator
}

// NewPaginatorFromParams takes an interface of type `PaginationParams`,
// the `url.Values` type works great with this interface, and returns
// a new `Paginator` based on the params or `PaginatorPageKey` and
// `PaginatorLimitKey`. Defaults are `1` for the page and
// PaginatorLimitDefault for the per page value.
func NewPaginatorFromParams(queryParams PaginatorParams) *Pagination {
	// Unmarshal the query params into a struct
	paginator := NewPaginatorWithDefaults()

	if limitStr := queryParams.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			if limit < 1 {
				limit = paginatorLimitDefault
			}

			paginator.Limit = limit
		}
	}

	if pageStr := queryParams.Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			if page < 1 {
				page = paginatorPageDefault
			}

			paginator.Page = page
			paginator.Offset = (page - 1) * paginator.Limit
		}
	}

	if sort := queryParams.Get("sort"); sort != "" {
		paginator.Sort = sort
	}

	return paginator
}

// Save - calculates and set pagination params based on records returned from the database.
func (p *Pagination) Save(count, responseLength int) {
	p.TotalEntriesSize = count
	p.CurrentEntriesSize = responseLength
	p.TotalPages = p.TotalEntriesSize / p.Limit

	if p.TotalEntriesSize%p.Limit > 0 {
		p.TotalPages++
	}
}

func (p *Pagination) GetLimit() int {
	return p.Limit
}

func (p *Pagination) GetOffset() int {
	return p.Offset
}

// Order returns ordering string.
func (p *Pagination) Order(defaultSortKey, defaultDir string) (key, direction string) {
	key, direction = parseSortField(p.Sort, nil)

	if direction == "" {
		direction = defaultDir
	}

	// direction must have a value either desc or asc (default).
	if !isValidDirection(direction) {
		direction = paginatorDefaultOrderDir
	}

	if key == "" {
		key = defaultSortKey
	}

	// Use the key and direction or default values to build the sort query
	return
}

// parseSortField parses the sort field into key and direction.
func parseSortField(sort string, keyToDBColumnMap map[string]string) (key, direction string) {
	if strings.Contains(sort, ":") {
		parts := strings.Split(sort, ":")
		key, direction = parts[0], parts[1]
	} else {
		key, direction = sort, ""
	}

	if keyToDBColumnMap != nil {
		// Map key to db column
		if dbColumnName, ok := keyToDBColumnMap[key]; ok {
			key = dbColumnName
		}
	}

	return
}

// isValidDirection checks if the given direction is valid.
func isValidDirection(direction string) bool {
	return strings.EqualFold(direction, "asc") || strings.EqualFold(direction, "desc")
}
