package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/adinovcina/golang-setup/tools/logger"
)

const (
	specificFileNameLength = 5
)

// FileNameRegex - Regex matches the following pattern:
// [version name bigint]_[Title].[up/down].sql
// Example:
//
//	123_name.up.sql
//	123_name.down.sql
var FileNameRegex = regexp.MustCompile(`^(\d+)_(.*)\.(down|up)\.(.*)$`)

// Migrate will start process of executing migration scripts against MySQL database server.
func (db *DB) Migrate(scriptsDirectory, scriptsEnvSpecificDirectory string) error {
	if err := db.db.Ping(); err != nil {
		return err
	}

	if err := MigrationSchema(db.db); err != nil {
		return err
	}

	fileNames, err := FilesInDirectory(scriptsDirectory, scriptsEnvSpecificDirectory)
	if err != nil {
		return fmt.Errorf("unable to read list of files from directory: %w", err)
	}

	filesExecutedList, err := FilesExecuted(db.db)
	if err != nil {
		return err
	}

	if len(fileNames) == 0 {
		logger.Info().Msg("# No Migration files found")
		return nil
	}

	logger.Info().Msgf("# Migration files found: %v", len(fileNames))

	keys := make([]string, 0, len(fileNames))

	for n := range fileNames {
		keys = append(keys, n)
	}

	sort.Strings(keys)

	filesExecuted := 0

	for _, name := range keys {
		if _, executed := filesExecutedList[name]; executed {
			continue
		}

		if err := executeMigrationFile(db.db, name, fileNames[name]); err != nil {
			return err
		}

		filesExecuted++
	}

	logger.Info().Msgf("# Migration finished: %v migration file(s) executed", filesExecuted)

	return nil
}

func executeMigrationFile(db *sql.DB, name, fileNameWithPath string) error {
	match := FileNameRegex.FindStringSubmatch(name)
	if len(match) != specificFileNameLength {
		return fmt.Errorf("error with parsing a file name %s with match params: %v", name, match)
	}

	fileBytes, err := os.ReadFile(fileNameWithPath)
	if err != nil {
		return err
	}

	queries := Parse(string(fileBytes))

	for _, statement := range queries {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("failed to execute statement: %v with error: %w", statement, err)
		}
	}

	logger.Info().Msgf("# Migration file executed: %v", name)

	return FileExecuted(match[0], match[1], match[2], db)
}
