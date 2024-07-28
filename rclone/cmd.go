package rclone

import (
	"encoding/json"
	"log"
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

func rcCopyDir(src, dest string) bool {
	args := append([]string{"copy", src, dest}, moveArgs...)
	lArgs := append(args, largeFileArgs...)
	if !rcExecCmd(lArgs...) {
		return false
	}
	sArgs := append(args, smallFileArgs...)
	return rcExecCmd(sArgs...)
}

func rcCopyFile(src, dest string) bool {
	args := append([]string{"copyto", src, dest, "--transfers=1", "--checkers=2"}, moveArgs...)
	return rcExecCmd(args...)
}

func rcExecCmd(args ...string) bool {
	cmdArgs := append([]string{"--quiet"}, args...)
	cmd := exec.Command("rclone", cmdArgs...)
	cmd.Env = cmdEnv

	var lastError []byte
COMMAND:
	for {
		var err error
		if lastError, err = cmd.CombinedOutput(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				switch exitError.ExitCode() {
				case exitcode.Success, exitcode.NoFilesTransferred:
					lastError = nil
				case exitcode.RetryError:
					lastError = nil
					continue COMMAND
				}
			}
		} else {
			lastError = nil
		}
		break
	}
	if lastError != nil {
		log.Print(string(lastError))
	}
	return lastError == nil
}
