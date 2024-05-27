package paging

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/adinovcina/golang-setup/tools/logger"
	"github.com/adinovcina/golang-setup/tools/utils"
)

const (
	prevCursorDirection = "previous"
	nextCursorDirection = "next"

	ascendingOrder  = "asc"
	descendingOrder = "desc"
)

type PaginatorCursor interface {
	Paginate(firstElem, lastElem any, resultLength int)
	Order(defaultSortKey, defaultDir string, keyToDBColumnMap map[string]string) (key, direction string)
	FormatLimit(limit int) string
	BuildWhereClause(where []string, args []interface{}, key, direction string) ([]string, []interface{}, error)

	GetLimit() int
	GetCursor() *string
	GetCursorDirection() string
	GetOrderByCursorDirection(direction string) string
}

// PaginatorCursorParams is a parameters provider interface to get the pagination params from.
type PaginatorCursorParams interface {
	Get(key string) string
}

type PaginationCursor struct {
	Cursor *Cursor `json:"cursor"`
	// Sort - Sort field in a format - sort=field:asc
	Sort string `json:"sort"`
	// Limit - A limit on the number of objects to be returned
	Limit int `json:"limit"`
}

type Cursor struct {
	PreviousCursor *string `json:"previousCursor"`
	NextCursor     *string `json:"nextCursor"`
	Cursor         string  `json:"-"`
	Direction      string  `json:"-"`
	HasNext        bool    `json:"hasNext"`
}

func NewPaginatorCursorWithDefaults() *PaginationCursor {
	return newPaginatorCursor(paginatorLimitDefault, nextCursorDirection)
}

// NewPaginatorCursorFromParams takes an interface of type `PaginationCursorParams`,
// the `url.Values` type works great with this interface, and returns a new `PaginatorCursor`.
func NewPaginatorCursorFromParams(queryParams PaginatorCursorParams) *PaginationCursor {
	// Unmarshal the query params into a struct
	paginator := NewPaginatorCursorWithDefaults()

	if limitStr := queryParams.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			if limit < 1 {
				limit = paginatorLimitDefault
			}

			paginator.Limit = limit
		}
	}

	if cursor := queryParams.Get("cursor"); cursor != "" {
		paginator.Cursor.Cursor = cursor
	}

	if direction := queryParams.Get("direction"); direction != "" {
		paginator.Cursor.Direction = direction
	}

	if sort := queryParams.Get("sort"); sort != "" {
		paginator.Sort = sort
	}

	return paginator
}

// newPaginatorCursor creates a new instance of PaginationCursor with the specified limit and direction.
func newPaginatorCursor(limit int, direction string) *PaginationCursor {
	// Ensure the provided limit is valid; if not, use the default limit
	if limit < 1 {
		limit = paginatorLimitDefault
	}

	// Ensure the direction is specified; if not, use the default direction
	if direction == "" {
		direction = nextCursorDirection
	}

	paginator := &PaginationCursor{
		Limit: limit,
		Cursor: &Cursor{
			Direction: direction,
		},
	}

	return paginator
}

// GetLimit returns paginator limit value.
func (p *PaginationCursor) GetLimit() int {
	return p.Limit
}

// GetCursor decodes and retrieves the cursor string from the pagination cursor.
// It returns a pointer to the decoded cursor string.
func (p *PaginationCursor) GetCursor() *string {
	// If the cursor is empty, return nil
	if p.Cursor.Cursor == "" {
		return nil
	}

	// Decode the cursor string from base64
	decodedCursor, err := decodeCursor(p.Cursor.Cursor)
	if err != nil {
		// Log an error if decoding fails
		logger.Error().Err(err).Msgf("failed to decode cursor %v", p.Cursor.Cursor)
		return nil
	}

	// Return a pointer to the decoded cursor string
	return utils.Ptr(decodedCursor.(string))
}

// GetOrderByCursorDirection returns the order direction based on the cursor direction and the requested direction.
func (p *PaginationCursor) GetOrderByCursorDirection(direction string) string {
	// Determine the cursor direction
	switch p.Cursor.Direction {
	case nextCursorDirection:
		// For "next" cursor, return the requested direction
		if direction == ascendingOrder || direction == descendingOrder {
			return direction
		}
	case prevCursorDirection:
		// For "previous" cursor, reverse the requested direction
		if direction == ascendingOrder {
			return descendingOrder
		} else if direction == descendingOrder {
			return ascendingOrder
		}
	}

	// Default to ascending order if direction is not recognized or cursor direction is unknown
	return ascendingOrder
}

func (p *PaginationCursor) GetCursorDirection() string {
	return p.Cursor.Direction
}

// Order returns ordering direction string.
func (p *PaginationCursor) Order(defaultSortKey, defaultDir string, keyToDBColumnMap map[string]string) (key, direction string) {
	key, direction = parseSortField(p.Sort, keyToDBColumnMap)

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

// BuildWhereClause builds a WHERE clause for pagination based on the cursor and sort direction.
// It takes a slice of strings representing the WHERE clause conditions, a slice of interface{}
// representing the arguments for prepared statements, the key to compare against, and the direction
// of pagination.
func (p *PaginationCursor) BuildWhereClause(
	where []string,
	args []interface{},
	key, direction string,
) ([]string, []interface{}, error) {
	// Retrieve current cursor and its direction
	cursor := p.GetCursor()
	cursorDirection := p.GetCursorDirection()

	if cursor == nil {
		return where, args, nil
	}

	// Construct WHERE clause based on cursor direction and sort direction
	switch {
	case cursorDirection == nextCursorDirection && direction == ascendingOrder:
		where, args = append(where, key+" > COALESCE(?, 0)"), append(args, *cursor)
	case cursorDirection == prevCursorDirection && direction == ascendingOrder:
		where, args = append(where, key+" < COALESCE(?, 0)"), append(args, *cursor)
	case cursorDirection == nextCursorDirection && direction == descendingOrder:
		// find max value for a key type
		maxValueForKey, err := findMax(*cursor)
		if err != nil {
			return nil, nil, fmt.Errorf("find max for %s: %w", *cursor, err)
		}

		where, args = append(where, fmt.Sprintf("%s < COALESCE(?, %s)", key, maxValueForKey)), append(args, *cursor)
	default:
		where, args = append(where, key+" > COALESCE(?, 0)"), append(args, *cursor)
	}

	return where, args, nil
}

// FormatLimit returns a SQL string for a given limit & offset.
// Clauses are only added if limit and/or offset are greater than zero.
func (p *PaginationCursor) FormatLimit(limit int) string {
	if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)
	}

	return ""
}

// Paginate saves pagination cursors based on the first and last elements of a result slice.
func (p *PaginationCursor) Paginate(firstElem, lastElem any, resultLength int) {
	// Determine cursor direction
	if p.GetCursorDirection() == prevCursorDirection {
		// Set nextCursor based on current cursor if it exists
		if p.GetCursor() != nil {
			p.Cursor.NextCursor = encodeCursor(lastElem)
		}

		// Set previousCursor based on the length of the result slice
		if resultLength > 0 && resultLength == p.GetLimit()+1 {
			p.Cursor.PreviousCursor = encodeCursor(firstElem)
		} else {
			p.Cursor.PreviousCursor = nil
		}
	} else {
		// Set previousCursor based on current cursor if it exists
		if p.GetCursor() != nil && resultLength > 0 {
			p.Cursor.PreviousCursor = encodeCursor(firstElem)
		}

		// Set nextCursor based on the length of the result slice
		if resultLength > p.GetLimit() && lastElem != nil {
			p.Cursor.NextCursor = encodeCursor(lastElem)
		}
	}

	// Set HasNext to true if nextCursor exists
	p.Cursor.HasNext = p.Cursor.NextCursor != nil
}

func findMax(input string) (string, error) {
	// Try parsing the input as an integer
	_, err := strconv.ParseInt(input, 10, 64)
	if err == nil {
		// Input is an integer, return the maximum possible positive integer value
		return strconv.FormatInt(int64(^uint64(0)>>1), 10), nil
	}

	// Try parsing the input as a date
	_, err = time.Parse("2006-01-02", input)
	if err == nil {
		// Input is a date, return the maximum possible date value
		maxDate := time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC)
		return maxDate.Format("2006-01-02"), nil
	}

	return "", fmt.Errorf("unsupported input: %s", input)
}

// encodeCursor encodes a cursor using base64 encoding.
func encodeCursor(elem any) *string {
	if elem == nil {
		return nil
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(elem)))

	return utils.Ptr(encoded)
}

// decodeCursor decodes a cursor previously encoded using base64.
func decodeCursor(encodedCursor string) (any, error) {
	if encodedCursor == "" {
		return nil, errors.New("encoded cursor must not be empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return nil, err
	}

	return string(decoded), nil
}
