package models

type PubKeyRequest struct {
	ClientVersion string `json:"clientVersion"`
	ConnectionId  string `json:"connectionId"`
	RemoteAddress string `json:"remoteAddress"`
	Username      string `json:"username"`
	PublicKey     string `json:"publicKey"`
}
