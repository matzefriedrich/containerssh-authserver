package models

import (
	"go.containerssh.io/containerssh/config"
)

type ConfigDockerResponse struct {
	AuthenticatedUsername string           `json:"authenticatedUsername"`
	ClientVersion         string           `json:"clientVersion"`
	Config                config.AppConfig `json:"config"`
	ConnectionId          string           `json:"connectionId"`
	RemoteAddress         string           `json:"remoteAddress"`
	Username              string           `json:"username"`
}
