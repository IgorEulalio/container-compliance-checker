package file

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// FileOptionCheck verifies that a file contains a specific option with optional expected value
type FileOptionCheck struct {
	FilePath      string `yaml:"file"`
	Option        string `yaml:"option"`
	ExpectedValue string `yaml:"value"`
}

// Name returns the name of this check
func (c *FileOptionCheck) Name() string {
	return "FileOption"
}

// Run executes the file option check
func (c *FileOptionCheck) Run() (bool, error) {
	file, err := os.Open(c.FilePath)
	if err != nil {
		return true, nil // If the file does not exist, we consider it a pass
	}
	defer file.Close()

	// Parse the file line by line
	scanner := bufio.NewScanner(file)
	optionFound := false
	valueMatches := c.ExpectedValue == "" // If no expected value is provided, we only check if the option exists

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and section headers
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "[") {
			continue
		}

		// Look for key=value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Check if this is the option we're looking for
			if key == c.Option {
				optionFound = true

				// If expected value is provided, check if it matches
				if c.ExpectedValue != "" && value == c.ExpectedValue {
					valueMatches = true
				}

				break // We found our option, no need to continue
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("error reading file %s: %w", c.FilePath, err)
	}

	if !optionFound {
		return false, fmt.Errorf("option %s not found in file %s", c.Option, c.FilePath)
	}

	if !valueMatches {
		return false, fmt.Errorf("option %s found but value does not match expected %s", c.Option, c.ExpectedValue)
	}

	return true, nil
}

// NewFileOptionCheck creates a new file option check with the given configuration
func NewFileOptionCheck(config map[string]interface{}) (*FileOptionCheck, error) {
	filePath, ok := config["file"]
	if !ok {
		return nil, fmt.Errorf("missing required 'file' field in configuration")
	}

	filePathStr, ok := filePath.(string)
	if !ok {
		return nil, fmt.Errorf("'file' field must be a string, got %T", filePath)
	}

	if filePathStr == "" {
		return nil, fmt.Errorf("'file' field cannot be empty")
	}

	option, ok := config["option"]
	if !ok {
		return nil, fmt.Errorf("missing required 'option' field in configuration")
	}

	optionStr, ok := option.(string)
	if !ok {
		return nil, fmt.Errorf("'option' field must be a string, got %T", option)
	}

	if optionStr == "" {
		return nil, fmt.Errorf("'option' field cannot be empty")
	}

	// Value is optional
	var valueStr string
	if value, ok := config["value"]; ok {
		if valueStr, ok = value.(string); !ok {
			return nil, fmt.Errorf("'value' field must be a string, got %T", value)
		}
	}

	return &FileOptionCheck{
		FilePath:      filePathStr,
		Option:        optionStr,
		ExpectedValue: valueStr,
	}, nil
}
