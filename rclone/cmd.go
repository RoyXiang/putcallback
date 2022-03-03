package rclone

import (
	"encoding/json"
	"errors"
	"os/exec"

	"github.com/rclone/rclone/lib/exitcode"
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

func rcMoveDir(src, dest string, transfers int, arg ...string) {
	args := []string{"move", src, dest}
	args = append(args, arg...)
	args = append(args, moveArgs...)
	rcExecCmd(transfers, args...)
}

func rcMoveFile(src, dest string, transfers int) {
	args := []string{"moveto", src, dest, "--transfers=1", "--checkers=2"}
	args = append(args, moveArgs...)
	rcExecCmd(transfers, args...)
}

func rcExecCmd(transfers int, args ...string) {
	for i := 0; i < transfers; i++ {
		transferQueue <- struct{}{}
	}
	defer func() {
		for i := 0; i < transfers; i++ {
			<-transferQueue
		}
	}()

	cmd := exec.Command("rclone", args...)
	var exitError *exec.ExitError
	for {
		if err := cmd.Run(); err != nil && errors.As(err, &exitError) && exitError.ExitCode() == exitcode.RetryError {
			continue
		}
		break
	}
}
