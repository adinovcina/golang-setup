package mysql

/*
	 This parse method will split query into statements where
	between each statement is --MYSQL_CUSTOM_STATEMENT_DELIMITER
	delimiter constant string
*/
import (
	"os"
	"path/filepath"
	"strings"
)

// Parse SQL query and return slice of statements.
func Parse(query string) []string {
	delimiter := "--MYSQL_CUSTOM_STATEMENT_DELIMITER"

	var result []string

	query = strings.TrimSpace(query)

	statements := strings.Split(query, delimiter)

	// Trim empty spaces
	for i, s := range statements {
		statements[i] = strings.TrimSpace(s)

		if len(statements[i]) > 0 {
			result = append(result, statements[i])
		}
	}

	return result
}

// FilesInDirectory retrieves the list of all files in the directory.
func FilesInDirectory(directory, envSpecificDirectory string) (map[string]string, error) {
	files := make(map[string]string)

	// Normalize directory paths to use forward slashes
	directory = filepath.ToSlash(directory)
	envSpecificDirectory = filepath.ToSlash(envSpecificDirectory)

	// Calculate absolute paths for comparison
	absDirectory := filepath.Clean(directory)
	absEnvSpecificDirectory := filepath.Join(absDirectory, envSpecificDirectory)

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Normalize path to use forward slashes
			path = filepath.ToSlash(path)

			// Check if the file is directly under the main directory or the environment specific directory
			if filepath.Dir(path) == absDirectory || filepath.Dir(path) == absEnvSpecificDirectory {
				files[info.Name()] = path
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
