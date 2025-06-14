package file

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FilePermissionCheck verifies that a file has the expected permissions
type FilePermissionCheck struct {
	FilePath    string `yaml:"file_path"`
	Permissions string `yaml:"permissions"`
}

// Name returns the name of this check
func (c *FilePermissionCheck) Name() string {
	return "FilePermission"
}

// Run executes the file permission check
func (c *FilePermissionCheck) Run() (bool, error) {
	fileInfo, err := os.Stat(c.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("file %s does not exist", c.FilePath)
		}
		return false, fmt.Errorf("error accessing file %s: %w", c.FilePath, err)
	}

	actualMode := fileInfo.Mode().Perm()

	// Parse expected permissions (from octal string like "0644")
	expectedMode, err := parsePermissions(c.Permissions)
	if err != nil {
		return false, err
	}

	if actualMode != expectedMode {
		return false, fmt.Errorf("file %s has permissions %o, expected %o",
			c.FilePath, actualMode, expectedMode)
	}

	return true, nil
}

// parsePermissions converts a permission string to os.FileMode
func parsePermissions(perms string) (os.FileMode, error) {
	// Remove leading "0" if present
	perms = strings.TrimPrefix(perms, "0")

	// Parse as octal
	modeInt, err := strconv.ParseUint(perms, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid permission format '%s': %w", perms, err)
	}

	return os.FileMode(modeInt), nil
}

// NewFilePermissionCheck creates a new file permission check with the given configuration
func NewFilePermissionCheck(config map[string]interface{}) (*FilePermissionCheck, error) {
	filePath, ok := config["file_path"]
	if !ok {
		return nil, fmt.Errorf("missing required 'file_path' field in configuration")
	}

	filePathStr, ok := filePath.(string)
	if !ok {
		return nil, fmt.Errorf("'file_path' field must be a string, got %T", filePath)
	}

	if filePathStr == "" {
		return nil, fmt.Errorf("'file_path' field cannot be empty")
	}

	permissions, ok := config["permissions"]
	if !ok {
		return nil, fmt.Errorf("missing required 'permissions' field in configuration")
	}

	permissionsStr, ok := permissions.(string)
	if !ok {
		return nil, fmt.Errorf("'permissions' field must be a string, got %T", permissions)
	}

	if permissionsStr == "" {
		return nil, fmt.Errorf("'permissions' field cannot be empty")
	}

	return &FilePermissionCheck{
		FilePath:    filePathStr,
		Permissions: permissionsStr,
	}, nil
}
