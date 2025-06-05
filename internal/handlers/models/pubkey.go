package models

type PubKeyRequest struct {
	ClientVersion string `json:"clientVersion"`
	ConnectionId  string `json:"connectionId"`
	RemoteAddress string `json:"remoteAddress"`
	Username      string `json:"username"`
	PublicKey     string `json:"publicKey"`
}

type PubKeyResponse struct {
	AuthenticatedUsername string `json:"authenticatedUsername"`
	ClientVersion         string `json:"clientVersion"`
	ConnectionId          string `json:"connectionId"`
	RemoteAddress         string `json:"remoteAddress"`
	Success               bool   `json:"success"`
	Username              string `json:"username"`
}
