package common

import "fmt"

// AppVersion is the application version taken from APP_VERSION in common.sh.
// GitRef is the git working copy HEAD the binary was built from.
// Both are set at build time via
// -ldflags "-X rpm-signing/common.AppVersion=... -X rpm-signing/common.GitRef=..."
// and keep their defaults when the binary is built without those flags.
var (
	AppVersion = "unknown"
	GitRef     = "unknown"
)

// VersionString returns the application version as "<version> (<git ref>)".
func VersionString() string {
	return fmt.Sprintf("%s (%s)", AppVersion, GitRef)
}
