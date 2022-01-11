package rclone

type RemoteConfig struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type RemotePutIoToken struct {
	AccessToken string `json:"access_token"`
	Expiry      string `json:"expiry"`
}

type EpisodeInfo struct {
	Group   string
	Show    string
	Season  int
	Episode string
	Version int
	Extra   string
}
