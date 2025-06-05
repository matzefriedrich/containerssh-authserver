package internal

var (
	// Version indicates the current version of the application, typically set at build time or defaulting to an empty string.
	Version = ""

	// CommitSha represents the commit SHA of the current build, defaulting to an empty string if not provided.
	CommitSha = ""

	// ReleaseDate holds the release date information as a string.
	ReleaseDate = ""

	// ReleaseName specifies the name of the release, used to identify the application version or deployment.
	ReleaseName = "authserver"
)
