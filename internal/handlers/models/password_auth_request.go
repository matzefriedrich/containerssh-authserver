package models

type PasswordRequest struct {
	Username       string `json:"username"`
	RemoteAddress  string `json:"remoteAddress"`
	ConnectionId   string `json:"connectionId"`
	ClientVersion  string `json:"clientVersion"`
	PasswordBase64 string `json:"passwordBase64"`
}
