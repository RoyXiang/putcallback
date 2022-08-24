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

func rcMoveDir(src, dest string) {
	args := append([]string{"move", src, dest}, moveArgs...)

	lArgs := append(args, largeFileArgs...)
	rcExecCmd(largeFileTransfers*2, lArgs...)

	sArgs := append(args, smallFileArgs...)
	rcExecCmd(smallFileTransfers, sArgs...)
}

func rcMoveFile(src, dest string, filesize int64) {
	args := append([]string{"moveto", src, dest, "--transfers=1", "--checkers=2"}, moveArgs...)
	if filesize < multiThreadCutoff {
		rcExecCmd(1, args...)
	} else {
		rcExecCmd(2, args...)
	}
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
	cmd.Env = cmdEnv

	var exitError *exec.ExitError
	for {
		if err := cmd.Run(); err != nil && errors.As(err, &exitError) && exitError.ExitCode() == exitcode.RetryError {
			continue
		}
		break
	}
}
