package rclone

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/RoyXiang/putcallback/putio"
)

var (
	renamingStyle string

	cmdEnv     []string
	fileChan   chan string
	folderChan chan string
	mu         sync.Mutex
	wg         sync.WaitGroup

	Put *putio.Put
)

func init() {
	rcEnv := []string{
		"RCLONE_DELETE_EMPTY_SRC_DIRS=true",
		"RCLONE_NO_TRAVERSE=true",
		"RCLONE_USE_MMAP=true",
	}
	cmdEnv = append(os.Environ(), rcEnv...)

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

	fileChan = make(chan string, 1)
	folderChan = make(chan string, Put.MaxTransfers)
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
