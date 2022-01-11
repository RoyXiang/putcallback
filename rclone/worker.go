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
	reFilename = regexp.MustCompile(`^(\[.+?])[\[ ](.+?)[] ]?-?[\[ ](?:EP)?([0-9]{1,3})(?:[vV]([0-9]))?[] ]((?:\[?END[] ])?[\[(].*)$`)
	reSeason   = regexp.MustCompile(`^S?([0-9]+)$`)
	reNumber   = regexp.MustCompile(`(\b|[ -])[0-9]+(\b|[ -])`)
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

	elems := strings.FieldsFunc(matches[2], func(r rune) bool {
		return r == ' ' || r == '[' || r == ']'
	})
	if matches := reSeason.FindStringSubmatch(elems[len(elems)-1]); matches != nil {
		season, _ := strconv.Atoi(matches[1])
		if season >= 100 {
			elems = elems[:len(elems)-1]
		} else {
			elems[len(elems)-1] = romanLib.ToRoman(season)
		}
	}
	i := 0
	for _, elem := range elems {
		elem = reNumber.ReplaceAllString(elem, "")
		elem = strings.Trim(elem, "-")
		if elem != "" && elem != "()" {
			elems[i] = elem
			i++
		}
	}
	name := strings.Join(elems[:i], " ")

	var episode string
	if matches[4] == "" {
		episode = matches[3]
	} else {
		episode = fmt.Sprintf("%sv%s", matches[3], matches[4])
	}
	return fmt.Sprintf("%s %s - %s %s", matches[1], name, episode, matches[5])
}
