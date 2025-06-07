package models

import (
	"github.com/matzefriedrich/containerssh-authserver/internal/shims"
)

type ConfigDockerResponse struct {
	AuthenticatedUsername string              `json:"authenticatedUsername"`
	ClientVersion         string              `json:"clientVersion"`
	Config                shims.AppConfigShim `json:"config" comment:"This is a shim type that has been added to make the config compatible with the latest Docker SDK."`
	ConnectionId          string              `json:"connectionId"`
	RemoteAddress         string              `json:"remoteAddress"`
	Username              string              `json:"username"`
}
