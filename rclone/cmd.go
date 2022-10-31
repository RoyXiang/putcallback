package rclone

import (
	"encoding/json"
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

func rcMoveDir(src, dest string) bool {
	args := append([]string{"move", src, dest}, moveArgs...)

	lArgs := append(args, largeFileArgs...)
	if !rcExecCmd(largeFileTransfers*2, lArgs...) {
		return false
	}

	sArgs := append(args, smallFileArgs...)
	return rcExecCmd(smallFileTransfers, sArgs...)
}

func rcMoveFile(src, dest string, filesize int64) bool {
	args := append([]string{"moveto", src, dest, "--transfers=1", "--checkers=2"}, moveArgs...)
	if filesize < multiThreadCutoff {
		return rcExecCmd(1, args...)
	}
	return rcExecCmd(2, args...)
}

func rcExecCmd(transfers int, args ...string) bool {
	for i := 0; i < transfers; i++ {
		transferQueue <- struct{}{}
	}
	defer func() {
		for i := 0; i < transfers; i++ {
			<-transferQueue
		}
	}()

	cmd := exec.Command("rclone", args...)
	cmd.Env = cmdEnv

	shouldRetry, hasErrors := true, false
	for shouldRetry {
		shouldRetry, hasErrors = false, false
		if err := cmd.Run(); err != nil {
			hasErrors = true
			if exitError, ok := err.(*exec.ExitError); ok {
				switch exitError.ExitCode() {
				case exitcode.Success, exitcode.NoFilesTransferred:
					hasErrors = false
				case exitcode.RetryError:
					shouldRetry = true
				}
			}
		}
	}
	return !hasErrors
}
