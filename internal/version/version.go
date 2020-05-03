package version

// injected values by ldflags
var (
	version string // $ git describe --tags
)

// Version returns the version.
// ex. 0.0.0-g-abababab
func Version() string {
	if len(version) == 0 {
		return "UNKNOWN_VERSION"
	}
	return version
}
