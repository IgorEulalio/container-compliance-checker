package packagecheck

import (
	"fmt"
	"os/exec"
)

// PackageNotInstalledCheck verifies that a specific package is not installed
type PackageNotInstalledCheck struct {
	Package string `yaml:"package"`
}

// CheckName returns the name of this check
func (c *PackageNotInstalledCheck) Name() string {
	return "PackageNotInstalled"
}

// PerformCheck executes the package installation check
func (c *PackageNotInstalledCheck) Run() (bool, error) {

	// Try different package managers
	pkgManagers := []struct {
		name    string
		command string
		args    []string
		success int // Expected exit code for "package not found"
	}{
		{"dpkg", "dpkg", []string{"-l", c.Package}, 1},                           // Debian/Ubuntu
		{"rpm", "rpm", []string{"-q", c.Package}, 1},                             // RHEL/CentOS/Fedora
		{"apk", "apk", []string{"info", "-e", c.Package}, 1},                     // Alpine
		{"pacman", "pacman", []string{"-Q", c.Package}, 1},                       // Arch
		{"zypper", "zypper", []string{"se", "--installed-only", c.Package}, 104}, // SUSE
	}

	for _, pm := range pkgManagers {
		if pmExists(pm.name) {
			cmd := exec.Command(pm.command, pm.args...)
			err := cmd.Run()

			if err != nil {
				// If exit status matches expected code for "not installed", check passes
				if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == pm.success {
					return true, nil
				}
				// Other errors could indicate problems running the command
				if exitErr, ok := err.(*exec.ExitError); ok {
					return false, fmt.Errorf("package check failed with exit code %d", exitErr.ExitCode())
				}
				return false, fmt.Errorf("failed to check package: %w", err)
			}

			// No error means package is installed, which fails the check
			return false, nil
		}
	}

	return false, fmt.Errorf("no supported package manager found")
}

// pmExists checks if a package manager exists in the system
func pmExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// NewPackageNotInstalledCheck creates a new package check with the given configuration
func NewPackageNotInstalledCheck(config map[string]interface{}) (*PackageNotInstalledCheck, error) {
	packageVal, ok := config["package"]
	if !ok {
		return nil, fmt.Errorf("missing required 'package' field in configuration")
	}

	packageName, ok := packageVal.(string)
	if !ok {
		return nil, fmt.Errorf("'package' field must be a string, got %T", packageVal)
	}

	if packageName == "" {
		return nil, fmt.Errorf("'package' field cannot be empty")
	}

	return &PackageNotInstalledCheck{
		Package: packageName,
	}, nil
}
