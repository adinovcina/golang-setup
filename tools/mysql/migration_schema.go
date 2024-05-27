package mysql

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

const (
	errorMySQLTableNotExists = 1146
)

// MigrationSchema validates if migration table exists if not try to add one.
func MigrationSchema(db *sql.DB) error {
	// Check if table exists
	err := migrationTableExists(db)

	// Get my sql error to check if table exists. Should return code 1146
	var mysqlError *mysql.MySQLError
	if !errors.As(err, &mysqlError) {
		// Unable  to convert sql error
		return err
	}

	// Table does not exists
	if mysqlError.Number == errorMySQLTableNotExists {
		// Create table
		createTableQuery := `CREATE TABLE IF NOT EXISTS migrations(
			id SERIAL, 
			file_name TEXT, 
			version TEXT, 
			title TEXT, 
			date_created DateTime NOT NULL DEFAULT CURRENT_TIMESTAMP, 
			PRIMARY KEY ( id )
			) ENGINE = InnoDB DEFAULT CHARSET=utf8`
		// Run query and do not handle results
		_, err := db.Exec(createTableQuery)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return err
		}

		err = migrationTableExists(db)
		if err != nil {
			return err
		}
	}

	return nil
}

// FilesExecuted retrieves the list of all executed files against database.
func FilesExecuted(db *sql.DB) (map[string]string, error) {
	query := "SELECT file_name FROM migrations"

	fileNames := make(map[string]string)

	rows, err := db.Query(query)
	// If there is an error but not empty row
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return fileNames, err
	}

	// This is valid sql error and means this file is not executed
	if errors.Is(err, sql.ErrNoRows) {
		return fileNames, nil
	}

	defer rows.Close()

	for rows.Next() {
		var fileName string

		err := rows.Scan(&fileName)
		if err != nil {
			return fileNames, err
		}

		fileNames[fileName] = fileName
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// File is executed
	return fileNames, nil
}

// FileExecuted will add file into migrations so it's not executed next time.
func FileExecuted(fileName, version, title string, db *sql.DB) error {
	query := "INSERT INTO migrations(file_name, version, title) VALUES(?, ?, ?)"

	result, err := db.Exec(query, fileName, version, title)
	if err != nil {
		return err
	}

	var id int64

	var e error
	if id, e = result.LastInsertId(); e != nil {
		return err
	}

	if id == 0 {
		return errors.New("migration not executed")
	}

	return nil
}

// migrationsExists validates if migrations exists in database migration schema.
func migrationTableExists(db *sql.DB) error {
	query := "SELECT 1 FROM migrations LIMIT 1"

	var count sql.NullFloat64

	err := db.QueryRow(query).Scan(&count)
	// There is err but not empty row
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	return nil
}
