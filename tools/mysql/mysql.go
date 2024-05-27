package mysql

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	// Using MySQL driver: https://github.com/go-sql-driver/mysql/
	_ "github.com/go-sql-driver/mysql" // Using blank import for the 'database/sql' module (needed internaly)
)

const (
	maxIdleConns    = 100
	maxOpenConns    = 100
	connMaxLifetime = 5
)

type DB struct {
	db *sql.DB

	ctx    context.Context // background context
	cancel func()          // cancel background context

	maxNumberOfRetries int
	retryBaseOffsetMs  int
}

// ConnectionData representing all connectivity data n eeded to establish a DB connection.
type ConnectionData struct {
	Address            string
	Port               string
	Name               string
	Username           string
	Password           string
	MaxNumberOfRetries string
	RetryBaseOffsetMs  string
}

// Init a DB connection suitable to the package's needs i.e. for the calling servise's DB.
func Init(connData *ConnectionData) (*DB, error) {
	db := new(DB)

	db.ctx, db.cancel = context.WithCancel(context.Background())

	if valErr := connData.Validate(); valErr != nil {
		return nil, valErr
	}

	var err error

	db.db, err = sql.Open("mysql", connData.ToConnectionStringWithoutDB())
	if err != nil {
		return nil, err
	}

	// Create the database if it doesn't exist
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", connData.Name)
	_, err = db.db.Exec(createDBQuery)
	if err != nil {
		return nil, err
	}

	// Close the connection
	defer db.db.Close()

	// Reopen the connection with the specified database name
	db.db, err = sql.Open("mysql", connData.ToConnectionString())
	if err != nil {
		return nil, err
	}

	db.db.SetConnMaxLifetime(time.Second * connMaxLifetime)
	db.db.SetMaxIdleConns(maxIdleConns)
	db.db.SetMaxOpenConns(maxOpenConns)

	err = db.db.Ping()
	if err != nil {
		return nil, err
	}

	// get retry values
	maxNumberOfRetries, retryBaseOffsetMs := connData.GetRetryValues()
	db.maxNumberOfRetries = maxNumberOfRetries
	db.retryBaseOffsetMs = retryBaseOffsetMs

	return db, nil
}

func (db *DB) GetDB() *sql.DB {
	return db.db
}

// Close closes the database connection.
func (db *DB) Close() error {
	// Cancel background context.
	db.cancel()

	// Close database.
	if db.db != nil {
		return db.db.Close()
	}

	return nil
}

// ToConnectionString convertes DB connectivity data into a mysql usable connection string.
func (cd *ConnectionData) ToConnectionString() string {
	return Concat(
		cd.Username,
		":",
		cd.Password,
		"@tcp(",
		cd.Address,
		":",
		cd.Port,
		")/",
		cd.Name,
		"?parseTime=true",
		"&charset=utf8mb4&collation=utf8mb4_unicode_ci",
	)
}

// ToConnectionStringWithoutDB convertes DB connectivity data into a mysql usable connection string without the database name.
func (cd *ConnectionData) ToConnectionStringWithoutDB() string {
	return Concat(
		cd.Username,
		":",
		cd.Password,
		"@tcp(",
		cd.Address,
		":",
		cd.Port,
		")/?parseTime=true",
		"&charset=utf8mb4&collation=utf8mb4_unicode_ci",
	)
}

// Validate the connectivity data.
func (cd *ConnectionData) Validate() error {
	if cd.Address == "" {
		return errors.New("invalid db connection data. no address provided")
	}

	if cd.Name == "" {
		return errors.New("invalid db connection data. no name provided")
	}

	if cd.Password == "" {
		return errors.New("invalid db connection data. no password provided")
	}

	if cd.Port == "" {
		return errors.New("invalid db connection data. no port provided")
	}

	if cd.Username == "" {
		return errors.New("invalid db connection data. no user name provided")
	}

	return nil
}

// GetRetryValues - Get values for retry execution.
func (cd *ConnectionData) GetRetryValues() (numberOfRetries, retryOffset int) {
	var num, offset int

	num, err := strconv.Atoi(cd.MaxNumberOfRetries)
	if err != nil || num == 0 {
		num = 3 // default value
	}

	offset, err = strconv.Atoi(cd.RetryBaseOffsetMs)
	if err != nil || offset == 0 {
		offset = 50 // default values
	}

	return num, offset
}

// Concat the provided strings.
func Concat(values ...string) string {
	var buffer bytes.Buffer

	for _, i := range values {
		buffer.WriteString(i)
	}

	return buffer.String()
}
