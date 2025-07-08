package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// FileRegexPresentCheck verifies that a regex pattern exists in at least one of the specified files
type FileRegexPresentCheck struct {
	files   []string `yaml:"files"`
	pattern string   `yaml:"pattern"`
}

// Name returns the name of this check
func (c *FileRegexPresentCheck) Name() string {
	return "FileRegexPresent"
}

// Run executes the file regex check
func (c *FileRegexPresentCheck) Run() (bool, error) {
	if len(c.files) == 0 {
		return false, fmt.Errorf("no files specified for regex check")
	}

	// Compile the regex pattern
	regex, err := regexp.Compile(c.pattern)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern '%s': %w", c.pattern, err)
	}

	// Track if we found the pattern in any file
	patternFound := false

	// Check each specified file or directory
	for _, filePath := range c.files {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue // Skip non-existent files
			}
			return false, fmt.Errorf("error accessing path %s: %w", filePath, err)
		}

		// If it's a directory, search all files in it
		if fileInfo.IsDir() {
			found, err := c.searchDirectory(filePath, regex)
			if err != nil {
				return false, err
			}
			if found {
				patternFound = true
				break
			}
		} else {
			// It's a regular file
			found, err := c.searchFile(filePath, regex)
			if err != nil {
				return false, err
			}
			if found {
				patternFound = true
				break
			}
		}
	}

	if !patternFound {
		return false, fmt.Errorf("pattern '%s' not found in any of the specified files", c.pattern)
	}

	return true, nil
}

// searchFile checks if the pattern exists in a single file
func (c *FileRegexPresentCheck) searchFile(filePath string, regex *regexp.Regexp) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if regex.MatchString(line) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return false, nil
}

// searchDirectory looks for the pattern in all files in a directory
func (c *FileRegexPresentCheck) searchDirectory(dirPath string, regex *regexp.Regexp) (bool, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return false, fmt.Errorf("error reading directory %s: %w", dirPath, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Skip subdirectories
		}

		fullPath := filepath.Join(dirPath, file.Name())
		found, err := c.searchFile(fullPath, regex)
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}

	return false, nil
}

// NewFileRegexPresentCheck creates a new file regex check with the given configuration
func NewFileRegexPresentCheck(config map[string]interface{}) (*FileRegexPresentCheck, error) {
	filesVal, ok := config["files"]
	if !ok {
		return nil, fmt.Errorf("missing required 'files' field in configuration")
	}

	// Extract files list
	var files []string
	filesArray, ok := filesVal.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'files' field must be an array, got %T", filesVal)
	}

	for _, f := range filesArray {
		fileStr, ok := f.(string)
		if !ok {
			return nil, fmt.Errorf("each file in 'files' must be a string, got %T", f)
		}
		files = append(files, fileStr)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("'files' field cannot be empty")
	}

	// Extract pattern
	patternVal, ok := config["pattern"]
	if !ok {
		return nil, fmt.Errorf("missing required 'pattern' field in configuration")
	}

	patternStr, ok := patternVal.(string)
	if !ok {
		return nil, fmt.Errorf("'pattern' field must be a string, got %T", patternVal)
	}

	if patternStr == "" {
		return nil, fmt.Errorf("'pattern' field cannot be empty")
	}

	return &FileRegexPresentCheck{
		files:   files,
		pattern: patternStr,
	}, nil
}
