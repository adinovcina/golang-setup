package mysqlstore

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/adinovcina/golang-setup/store"
	"github.com/adinovcina/golang-setup/tools/logger"
	"github.com/twinj/uuid"
)

// GetUserByID will retrieve the user filtered by ID
func (r *Repository) GetUserByID(id uuid.UUID) (*store.User, error) {
	query, err := r.db.Prepare("CALL GetUserByID(?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL GetUserByID(%d).", id)
		return nil, err
	}

	defer query.Close()

	user := new(store.User)

	err = query.QueryRow(id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Language, &user.Active,
			&user.Role, &user.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New(store.UserNotFound)
	}

	if err != nil {
		logger.Error().Err(err).Msgf("There was an error executing query: CALL GetUserByID(%d).", id)
		return nil, err
	}

	return user, nil
}

// Users will retrieve list of all users for specific admin based on table parameters and search requests.
func (r *Repository) GetUsers(filter *store.UserFilter) ([]*store.User, error) {
	limit := r.PaginatorCursor().GetLimit()
	// Get order key and direction values, with defaults specified
	key, direction := r.PaginatorCursor().Order("u.id", "desc", store.UsersKeyToColumnMap())

	var err error
	// Build WHERE clause. Each segment of the clause is AND-ed together.
	// Values are appended to args so we can avoid SQL injection.
	where, args := []string{"1 = 1"}, []interface{}{}

	// Build paginator where clause based on page direction
	where, args, err = r.PaginatorCursor().BuildWhereClause(where, args, key, direction)
	if err != nil {
		return nil, fmt.Errorf("build where clause: %w", err)
	}

	if v := filter.Active; v != nil {
		where, args = append(where, "u.active = ?"), append(args, *v)
	}
	if v := filter.Search; v != nil {
		// Instead of using CONCAT in SQL, construct the pattern string directly
		pattern := "%" + *v + "%"
		where, args = append(where, "u.name LIKE ?"), append(args, pattern)
	}

	// Determine default sorting.
	sortBy := fmt.Sprintf("%s %s",
		key, r.paginatorCursor.GetOrderByCursorDirection(direction))

	query := `SELECT u.id,
		u.name,
		u.email,
		u.phone,
		u.active,
		u.created_at,
		r.name
	FROM
		users u
	JOIN user_roles ur ON u.id = ur.user_id
	JOIN roles r ON r.id = ur.role_id
	WHERE `

	query += strings.Join(where, " AND ")
	query += ` ORDER BY ` + sortBy + `
	` + r.PaginatorCursor().FormatLimit(limit+1)

	// Limit is equal to count plus one, to fetch one more result than the count specified by the client.
	// The extra result isn’t returned in the result set, but we use the ID of the value as the next_cursor.

	// First we need to get the correct set of data,
	// and then get the correct set in the correct order (reverse the set)
	if r.PaginatorCursor().GetCursorDirection() == "previous" {
		sortByReverse := fmt.Sprintf("%s %s", key, direction)

		// For previous page, the main query will be subquery,
		// because we need to reverse the result set order again to
		// get the correct order
		subQuery := query

		query = `SELECT u.id,
					u.name,
					u.email,
					u.phone,
					u.active,
					u.created_at,
					r.name
				FROM(
					` + subQuery + `
				) AS bl
				ORDER BY ` + sortByReverse + ``
	}

	users := make([]*store.User, 0)

	// Prepare the query
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(args...)
	if errors.Is(err, sql.ErrNoRows) {
		return users, nil
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := new(store.User)

		err := rows.Scan(&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Active,
			&user.CreatedAt,
			&user.Role)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// get result length before slicing
	resultLength := len(users)

	if len(users) > 0 {
		if r.PaginatorCursor().GetCursorDirection() == "next" {
			if len(users) > limit {
				users = users[:limit]
			}
		} else {
			if len(users) > limit {
				users = users[1:]
			}
		}

		// Paginate
		r.PaginatorCursor().Paginate(users[0].ID, users[len(users)-1].ID, resultLength)
	}

	return users, nil
}

// GetUserByEmail will retrieve the user filtered by email
func (r *Repository) GetUserByEmail(email string) (*store.User, error) {
	query, err := r.db.Prepare("CALL GetUserByEmail(?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL GetUserByEmail(%s).", email)
		return nil, err
	}

	defer query.Close()

	user := new(store.User)

	err = query.QueryRow(email).
		Scan(&user.ID, &user.Name, &user.Email,
			&user.Password, &user.Active, &user.FailedLoginCount, &user.LoginBlockedUntil)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New(store.UserNotFound)
	}

	if err != nil {
		logger.Error().Err(err).Msgf("There was an error executing query: CALL GetUserByEmail(%s).", email)
		return nil, err
	}

	return user, nil
}

// GetUserByToken will retrieve the user filtered by temproary login token and info if token expired
func (r *Repository) GetUserByToken(token, tokenType string) (*store.User, error) {
	query, err := r.db.Prepare("CALL GetUserByToken(?,?)")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to prepare statement: CALL GetUserByToken(%v,%v).",
			token, tokenType)
		return nil, err
	}

	defer query.Close()

	user := new(store.User)

	err = query.QueryRow(token, tokenType).
		Scan(&user.Expired, &user.ID, &user.Name, &user.Email, &user.Active,
			&user.Role, &user.RoleID, &user.Language, &user.FailedLoginCount, &user.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New(store.UserNotFound)
	}

	if err != nil {
		logger.Error().Err(err).Msgf("There was an error executing query: GetUserByToken(%v,%v).",
			token, tokenType)
		return nil, err
	}

	return user, nil
}
