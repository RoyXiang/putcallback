package rclone

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/RoyXiang/putcallback/putio"
	"github.com/rclone/rclone/fs"
)

var (
	renamingStyle string

	moveArgs      []string
	largeFileArgs []string
	smallFileArgs []string

	fileChan   chan *putio.FileInfo
	folderChan chan *putio.FileInfo
	mu         sync.Mutex
	wg         sync.WaitGroup
	wgFolder   sync.WaitGroup

	Put *putio.Put
)

func init() {
	moveArgs = []string{
		"--check-first",
		"--no-traverse",
		"--use-mmap",
		"--drive-pacer-min-sleep=1ms",
	}
	rcGlobalConfig := fs.GetConfig(nil)
	largeFileArgs = []string{
		fmt.Sprintf("--transfers=%d", rcGlobalConfig.Transfers),
		fmt.Sprintf("--checkers=%d", rcGlobalConfig.Checkers),
		fmt.Sprintf("--min-size=%db", rcGlobalConfig.MultiThreadCutoff),
		"--delete-empty-src-dirs",
	}
	smallFileArgs = []string{
		fmt.Sprintf("--transfers=%d", rcGlobalConfig.Transfers*2),
		fmt.Sprintf("--checkers=%d", rcGlobalConfig.Checkers*2),
		fmt.Sprintf("--max-size=%db", rcGlobalConfig.MultiThreadCutoff-1),
		"--delete-empty-src-dirs",
	}

	styleInEnv := strings.ToLower(os.Getenv("RENAMING_STYLE"))
	if styleInEnv == RenamingStyleAnime {
		renamingStyle = RenamingStyleAnime
	} else if styleInEnv == RenamingStyleTv {
		renamingStyle = RenamingStyleTv
	} else {
		renamingStyle = RenamingStyleNone
	}

	accessToken := parseRCloneConfig()
	Put = putio.New(accessToken)

	fileChan = make(chan *putio.FileInfo, 1)
	folderChan = make(chan *putio.FileInfo, Put.MaxTransfers)
	wg.Add(2)
	go worker(fileChan)
	go moveFolder(folderChan)
}

func Stop() {
	close(fileChan)
	close(folderChan)
	wg.Wait()
}

func parseRCloneConfig() (accessToken string) {
	config := rcDumpConfig()
	if config == nil {
		log.Fatal("Please install rclone and configure it correctly")
	}
	if src, ok := config[RemoteSource]; !ok {
		log.Fatalf("Please configure %s as a rclone remote", RemoteSource)
	} else if src.Type != "putio" {
		log.Fatalf("Please ensure %s is a configuration of Put.io", RemoteSource)
	} else if src.Token != "" {
		var token RemotePutIoToken
		if err := json.Unmarshal([]byte(src.Token), &token); err != nil {
			log.Fatalf("Please ensure %s has a valid Put.io token", RemoteSource)
		}
		accessToken = token.AccessToken
	} else {
		log.Fatalf("Please configure %s correctly", RemoteSource)
	}
	if _, ok := config[RemoteDestination]; !ok {
		log.Fatalf("Please configure %s as a rclone remote", RemoteDestination)
	}
	return
}
