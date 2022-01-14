package rclone

import (
	"encoding/json"
	"errors"
	"os/exec"
)

func rcDumpConfig() map[string]RemoteConfig {
	output, err := exec.Command("rclone", "config", "dump").Output()
	if err != nil {
		return nil
	}
	var config map[string]RemoteConfig
	if err := json.Unmarshal(output, &config); err != nil {
		return nil
	}
	return config
}

func rcMoveDir(src, dest string, arg ...string) {
	defer wgFolder.Done()
	args := []string{"move", src, dest}
	args = append(args, arg...)
	rcExecCmd(args...)
}

func rcMoveFile(src, dest string, arg ...string) {
	args := []string{"moveto", src, dest}
	args = append(args, arg...)
	rcExecCmd(args...)
}

func rcExecCmd(args ...string) {
	cmd := exec.Command("rclone", args...)
	cmd.Env = cmdEnv

	var exitError *exec.ExitError
	for {
		if err := cmd.Run(); err != nil && errors.As(err, &exitError) && exitError.ExitCode() == ErrorTemporary {
			continue
		}
		break
	}
}
