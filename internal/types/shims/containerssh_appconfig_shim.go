// Package shims
// Note: This file contains shims for the AppConfigShim type (and related types) to make the configuration work with the latest github.com/docker/docker package.
package shims

type Backend string

const (
	BackendDocker Backend = "docker"
)

type AppConfigShim struct {
	Backend Backend          `json:"backend" default:"docker"`
	Docker  DockerConfigShim `json:"docker,omitempty"`
}
