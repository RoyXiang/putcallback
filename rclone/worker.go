package rclone

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/RoyXiang/putcallback/notification"
	"github.com/chonla/roman-number-go"
)

var (
	reFilename = regexp.MustCompile(`^(\[.+])[\[ ](.+?)[] ]?-?[\[ ](?:EP)?([0-9]+)(?:[vV]([0-9]{1,3}))?[] ]((\[?END[] ])?[\[(].*)$`)
	reSeason   = regexp.MustCompile(` S?(\d+)$`)
	romanLib   = roman.NewRoman()
)

func SendFileIdToWorker(fileId int64) {
	mu.Lock()
	defer mu.Unlock()

	name, isDir := Put.GetFileInfo(fileId)
	if name == "" {
		return
	}
	go Put.CleanupTransfers()
	if isDir {
		folderChan <- name
	} else {
		fileChan <- name
	}
}

func worker(fileChan <-chan string) {
	defer wg.Done()
	for filename := range fileChan {
		wg.Add(1)
		go moveFile(filename)
	}
}

func moveFolder(folderChan <-chan string) {
	defer wg.Done()
	for folder := range folderChan {
		log.Printf("Moving folder %s...", folder)

		src := fmt.Sprintf("%s:%s", RemoteSource, folder)
		dest := fmt.Sprintf("%s:%s", RemoteDestination, folder)
		rcMove(src, dest, "--transfers=20", "--checkers=30", "--max-size=250M")
		rcMove(src, dest, "--transfers=5", "--checkers=10", "--multi-thread-streams=10", "--min-size=250M", "--tpslimit=100", "--tpslimit-burst=100")
		rcRemoveDir(src)

		notification.Send(fmt.Sprintf("%s finished", folder))
	}
}

func moveFile(filename string) {
	defer wg.Done()

	log.Printf("Moving file %s...", filename)

	src := fmt.Sprintf("%s:%s", RemoteSource, filename)
	dest := fmt.Sprintf("%s:%s", RemoteDestination, RenameFile(filename))
	rcMoveTo(src, dest, "--transfers=1", "--checkers=2", "--multi-thread-streams=10", "--tpslimit=100", "--tpslimit-burst=100")

	notification.Send(fmt.Sprintf("%s finished", filename))
}

func RenameFile(filename string) string {
	matches := reFilename.FindStringSubmatch(filename)
	if matches == nil {
		return filename
	}
	mSeason := reSeason.FindStringSubmatch(matches[2])
	if mSeason != nil {
		season, _ := strconv.Atoi(mSeason[1])
		parts := strings.Split(matches[2], " ")
		parts[len(parts)-1] = romanLib.ToRoman(season)
		matches[2] = strings.Join(parts, " ")
	}
	var episode string
	if matches[4] == "" {
		episode = matches[3]
	} else {
		episode = fmt.Sprintf("%sv%s", matches[3], matches[4])
	}
	return fmt.Sprintf("%s %s - %s %s", matches[1], matches[2], episode, matches[5])
}
