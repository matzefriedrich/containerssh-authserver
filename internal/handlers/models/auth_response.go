package models

type AuthResponse struct {
	AuthenticatedUsername string `json:"authenticatedUsername"`
	ClientVersion         string `json:"clientVersion"`
	ConnectionId          string `json:"connectionId"`
	RemoteAddress         string `json:"remoteAddress"`
	Success               bool   `json:"success"`
	Username              string `json:"username"`
}
