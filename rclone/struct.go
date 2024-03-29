package rclone

import (
	"github.com/rclone/rclone/fs/fspath"
)

type Remote fspath.Parsed

type RemoteConfig struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type RemotePutIoToken struct {
	AccessToken string `json:"access_token"`
	Expiry      string `json:"expiry"`
}
