package rclone

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/RoyXiang/putcallback/putio"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/fspath"
	"golang.org/x/sync/semaphore"
)

var (
	remoteSrc  *Remote
	remoteDest *Remote

	Put *putio.Put

	renamingStyle       string
	delayBeforeTransfer time.Duration
	excludeFileTypes    []string

	cmdEnv        []string
	moveArgs      []string
	largeFileArgs []string
	smallFileArgs []string
	folderWeight  int64

	taskChan    chan *putio.FileInfo
	transferSem *semaphore.Weighted

	callbackMu sync.Mutex
	workerWg   sync.WaitGroup
)

func init() {
	rcGlobalConfig := fs.GetConfig(nil)
	argMultiThreadCutoff := int64(rcGlobalConfig.MultiThreadCutoff)
	argLargeFileTransfers := int64(rcGlobalConfig.Transfers)
	argSmallFileTransfers := argLargeFileTransfers * 2

	moveArgs = []string{
		"--check-first",
		"--no-traverse",
		"--use-mmap",
	}
	largeFileArgs = make([]string, 0, 4)
	smallFileArgs = make([]string, 0, 3)

	osEnv := os.Environ()
	maxTransfers := 0
	for _, env := range osEnv {
		pair := strings.SplitN(env, "=", 2)
		switch pair[0] {
		case "MAX_TRANSFERS":
			if maxTransfersInEnv, err := strconv.Atoi(pair[1]); err == nil {
				maxTransfers = maxTransfersInEnv
			}
		case "RENAMING_STYLE":
			styleInEnv := strings.ToLower(pair[1])
			switch styleInEnv {
			case RenamingStyleAnime, RenamingStyleTv:
				renamingStyle = styleInEnv
			default:
				renamingStyle = RenamingStyleNone
			}
		case "DELAY_BEFORE_TRANSFER":
			delayBeforeTransfer = 0
			if pair[1] != "" {
				if parsed, err := time.ParseDuration(pair[1]); err == nil {
					delayBeforeTransfer = parsed
				}
			}
		case "EXCLUDE_FILETYPES":
			excludeFileTypes = strings.FieldsFunc(pair[1], func(r rune) bool {
				return r == ',' || r == '.'
			})
			if len(excludeFileTypes) > 0 {
				filterArgs := fmt.Sprintf("--exclude=*.{%s}", strings.Join(excludeFileTypes, ","))
				largeFileArgs = append(largeFileArgs, filterArgs)
				smallFileArgs = append(smallFileArgs, filterArgs)
			}
		case "RCLONE_TRANSFERS":
			transfers, err := strconv.ParseInt(pair[1], 10, 64)
			if err != nil {
				break
			}
			if transfers > argLargeFileTransfers {
				argSmallFileTransfers = transfers * 2
			}
			argLargeFileTransfers = transfers
		default:
			if pair[0] == "HOME" || strings.HasPrefix(pair[0], "RCLONE_") {
				cmdEnv = append(cmdEnv, env)
			}
		}
	}

	accessToken := parseRCloneConfig()
	Put = putio.New(accessToken, maxTransfers)

	maxWeight := argLargeFileTransfers + 1
	folderWeight = maxWeight - 1
	largeFileArgs = append(
		largeFileArgs,
		fmt.Sprintf("--min-size=%db", argMultiThreadCutoff),
		fmt.Sprintf("--transfers=%d", argLargeFileTransfers),
		fmt.Sprintf("--checkers=%d", argLargeFileTransfers*2),
	)
	smallFileArgs = append(
		smallFileArgs,
		fmt.Sprintf("--transfers=%d", argSmallFileTransfers),
		fmt.Sprintf("--checkers=%d", argSmallFileTransfers*2),
	)

	taskChan = make(chan *putio.FileInfo, 1)
	transferSem = semaphore.NewWeighted(maxWeight)
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
