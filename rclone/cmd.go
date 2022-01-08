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

func rcMove(src, dest string, arg ...string) {
	args := []string{"move", src, dest}
	args = append(args, arg...)
	rcExecCmd(args...)
}

func rcMoveTo(src, dest string, arg ...string) {
	args := []string{"moveto", src, dest}
	args = append(args, arg...)
	rcExecCmd(args...)
}

func rcRemoveDir(dir string) {
	rcExecCmd("rmdir", dir)
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
