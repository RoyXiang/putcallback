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
	reFilename = regexp.MustCompile(`^(\[.+?])[\[ ](.+?)[] ]?-?[\[ ](E|EP|SP)?([0-9]{1,3}(?:\.[0-9])?)(?:[vV]([0-9]))?(?:\((OAD|OVA)\))?[] ]((?:\[?END[] ])?[\[(].*)$`)
	reSeason   = regexp.MustCompile(`^S?([0-9]+)$`)
	reDigits   = regexp.MustCompile(`(\b|-)[0-9]+(\b|-)`)
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

	var newFilename string
	if renamingStyle == RenamingStyleAnime {
		newFilename = RenameFileInAnimeStyle(filename)
	} else {
		newFilename = RenameFileInTvStyle(filename)
	}

	src := fmt.Sprintf("%s:%s", RemoteSource, filename)
	dest := fmt.Sprintf("%s:%s", RemoteDestination, newFilename)
	rcMoveTo(src, dest, "--transfers=1", "--checkers=2", "--multi-thread-streams=10", "--tpslimit=100", "--tpslimit-burst=100")

	notification.Send(fmt.Sprintf("%s finished", filename))
}

func ParseEpisodeInfo(filename string) *EpisodeInfo {
	matches := reFilename.FindStringSubmatch(filename)
	if matches == nil {
		return nil
	}

	info := &EpisodeInfo{
		Group:   matches[1],
		Season:  1,
		Episode: matches[4],
		Extra:   matches[7],
	}

	if matches[5] != "" {
		info.Version, _ = strconv.Atoi(matches[5])
	} else {
		info.Version = 1
	}

	elems := strings.FieldsFunc(matches[2], func(r rune) bool {
		return r == ' ' || r == '[' || r == ']'
	})
	lenElems := len(elems)
	if matches := reSeason.FindStringSubmatch(elems[lenElems-1]); matches != nil {
		season, _ := strconv.Atoi(matches[1])
		if lenElems >= 2 && strings.ToLower(elems[lenElems-2]) == "part" {
			elems[lenElems-1] = romanLib.ToRoman(season)
		} else if season < 100 {
			info.Season = season
			elems = elems[:lenElems-1]
		}
	}
	if matches[3] == "SP" || matches[6] != "" {
		info.Season = 0
	}
	info.Show = strings.Join(elems, " ")
	return info
}

func RenameFileInAnimeStyle(filename string) string {
	info := ParseEpisodeInfo(filename)
	if info == nil {
		return filename
	}

	elems := strings.Split(info.Show, " ")
	i := 0
	for _, elem := range elems {
		elem = reDigits.ReplaceAllString(elem, "")
		elem = strings.Trim(elem, "-")
		if elem != "" && elem != "()" {
			elems[i] = elem
			i++
		}
	}
	elems = elems[:i]
	if info.Season > 1 {
		elems = append(elems, romanLib.ToRoman(info.Season))
	}
	name := strings.Join(elems, " ")

	var prefix string
	if info.Season == 0 {
		prefix = "S"
	}
	var episode string
	if info.Version == 1 {
		episode = fmt.Sprintf("%s%s", prefix, info.Episode)
	} else {
		episode = fmt.Sprintf("%s%sv%d", prefix, info.Episode, info.Version)
	}
	return fmt.Sprintf("%s %s - %s %s", info.Group, name, episode, info.Extra)
}

func RenameFileInTvStyle(filename string) string {
	info := ParseEpisodeInfo(filename)
	if info == nil {
		return filename
	}
	return fmt.Sprintf("%s - S%02dE%s - %s%s", info.Show, info.Season, info.Episode, info.Group, info.Extra)
}
