package models

type ConfigRequest struct {
	AuthenticatedUsername string `json:"authenticatedUsername"`
	ClientVersion         string `json:"clientVersion"`
	ConnectionId          string `json:"connectionId"`
	RemoteAddr            string `json:"remoteAddress"`
	Username              string `json:"username"`
}
