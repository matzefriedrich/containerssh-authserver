package configuration

import (
	"strings"
)

type ApplicationInfo struct {
	ReleaseName       string
	ReleaseDate       string
	Version           string
	Commit            string
	ConfigurationPath string
	ConfigurationType string
}

func (a *ApplicationInfo) VersionString() string {

	buf := strings.Builder{}

	version := strings.TrimLeft(a.Version, "v")

	if len(version) > 0 {

		buf.WriteString(" v")
		buf.WriteString(version)

		if len(a.Commit) > 0 {
			buf.WriteString("#")
			buf.WriteString(a.Commit)
		}

		if len(a.ReleaseDate) > 0 {
			buf.WriteString(", ")
			buf.WriteString(a.ReleaseDate)
		}
	}

	return buf.String()
}
