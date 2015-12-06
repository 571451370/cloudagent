package qga

var (
	// Version string (git descrive --long)
	Version string
	// BuildTime build time
	BuildTime string
)

// GetVersion display current cloudagent version
func GetVersion() string {
	return Version
}
