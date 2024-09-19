package build

// These variables will be set at build time using ldflags
var (
	// Version is the version of the app with default value dev
	Version = "dev"

	// Time when the app was built
	Time string

	// User who build the app
	User string
)
