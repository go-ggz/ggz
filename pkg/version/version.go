package version

import (
	"fmt"
	"runtime"
)

var (
	// Version indicates development branch. Releases will be empty string.
	Version string
	// BuildDate is the ISO 8601 day drone was built.
	BuildDate string
)

// PrintCLIVersion print server info
func PrintCLIVersion() string {
	return fmt.Sprintf("version %s, built on %s, %s", Version, BuildDate, runtime.Version())
}
