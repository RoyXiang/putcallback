package rclone

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/RoyXiang/putcallback/putio"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/fspath"
)

var (
	remoteSrc  *Remote
	remoteDest *Remote

	Put *putio.Put

	renamingStyle       string
	delayBeforeTransfer time.Duration

	multiThreadCutoff  int64
	largeFileTransfers int
	smallFileTransfers int
	maxTransfers       int

	moveArgs      []string
	largeFileArgs []string
	smallFileArgs []string

	taskChan      chan *putio.FileInfo
	transferQueue chan struct{}

	callbackMu sync.Mutex
	folderMu   sync.Mutex
	workerWg   sync.WaitGroup
)

func init() {
	accessToken := parseRCloneConfig()
	Put = putio.New(accessToken)

	rcGlobalConfig := fs.GetConfig(nil)
	multiThreadCutoff = int64(rcGlobalConfig.MultiThreadCutoff)
	largeFileTransfers = rcGlobalConfig.Transfers
	smallFileTransfers = rcGlobalConfig.Transfers * 2
	maxTransfers = smallFileTransfers + 2

	moveArgs = []string{
		"--check-first",
		"--no-traverse",
		"--use-mmap",
		"--drive-pacer-min-sleep=1ms",
	}
	largeFileArgs = []string{
		fmt.Sprintf("--transfers=%d", largeFileTransfers),
		fmt.Sprintf("--checkers=%d", rcGlobalConfig.Checkers),
		fmt.Sprintf("--min-size=%db", multiThreadCutoff),
	}
	smallFileArgs = []string{
		fmt.Sprintf("--transfers=%d", smallFileTransfers),
		fmt.Sprintf("--checkers=%d", rcGlobalConfig.Checkers*2),
		fmt.Sprintf("--max-size=%db", multiThreadCutoff-1),
	}

	styleInEnv := strings.ToLower(os.Getenv("RENAMING_STYLE"))
	if styleInEnv == RenamingStyleAnime {
		renamingStyle = RenamingStyleAnime
	} else if styleInEnv == RenamingStyleTv {
		renamingStyle = RenamingStyleTv
	} else {
		renamingStyle = RenamingStyleNone
	}

	delayBeforeTransfer = 0
	delayInEnv := os.Getenv("DELAY_BEFORE_TRANSFER")
	if delayInEnv != "" {
		if parsed, err := time.ParseDuration(delayInEnv); err == nil {
			delayBeforeTransfer = parsed
		}
	}

	taskChan = make(chan *putio.FileInfo, 1)
	transferQueue = make(chan struct{}, maxTransfers)
}

func Start() {
	workerWg.Add(1)
	go worker()
}

func Stop() {
	close(taskChan)
	workerWg.Wait()
}

func parseRemote(env, defaultPath string) *Remote {
	path := os.Getenv(env)
	if path == "" {
		path = defaultPath
	}
	if parsed, err := fspath.Parse(path); err == nil {
		return (*Remote)(&parsed)
	}
	return nil
}

func parseRCloneConfig() (accessToken string) {
	remoteSrc = parseRemote("REMOTE_SRC", RemoteSource)
	if remoteSrc == nil || remoteSrc.Name == "" {
		log.Fatal("Invalid REMOTE_SRC value")
	}
	remoteDest = parseRemote("REMOTE_DEST", RemoteDestination)
	if remoteDest == nil {
		log.Fatal("Invalid REMOTE_DEST value")
	}

	config := rcDumpConfig()
	if config == nil {
		log.Fatal("Please install rclone and configure it correctly")
	}

	if src, ok := config[remoteSrc.Name]; !ok {
		log.Fatalf("Please configure REMOTE_SRC (%s) as a rclone remote", remoteSrc.Name)
	} else if src.Type != "putio" {
		log.Fatalf("Please ensure REMOTE_SRC (%s) is a configuration of Put.io", remoteSrc.Name)
	} else if src.Token != "" {
		var token RemotePutIoToken
		if err := json.Unmarshal([]byte(src.Token), &token); err != nil {
			log.Fatalf("Please ensure REMOTE_SRC (%s) has a valid Put.io token", remoteSrc.Name)
		}
		accessToken = token.AccessToken
	} else {
		log.Fatalf("Please configure REMOTE_SRC (%s) correctly", remoteSrc.Name)
	}

	if remoteDest.Name != "" {
		if _, ok := config[remoteDest.Name]; !ok {
			log.Fatalf("Please configure REMOTE_DEST (%s) as a rclone remote", remoteDest.Name)
		}
	}
	return
}
