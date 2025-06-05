package configuration

import "strings"

type ApplicationInfo struct {
	ReleaseName       string
	Version           string
	Commit            string
	ConfigurationPath string
	ConfigurationType string
}

func (a *ApplicationInfo) VersionString() string {

	buf := strings.Builder{}

	if len(a.ReleaseName) > 0 {
		buf.WriteString(a.ReleaseName)
		buf.WriteString(" ")
	}

	version := strings.TrimLeft(a.Version, "v")
	buf.WriteString("v")
	buf.WriteString(version)

	if len(a.Commit) > 0 {
		buf.WriteString("#")
		buf.WriteString(a.Commit)
	}

	return buf.String()
}
