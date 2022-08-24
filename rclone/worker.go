package rclone

import (
	"fmt"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/RoyXiang/putcallback/notification"
	"github.com/RoyXiang/putcallback/putio"
	"github.com/chonla/roman-number-go"
)

var (
	reFilename = regexp.MustCompile(`^(\[.+?])[\[ ]?(.+?)[] ]?-?[\[ ](E|EP|SP)?([0-9]{2,3}(?:\.[0-9])?)(?:[vV]([0-9]))?(?:\((.+)\))?[] ]((?:\[?END[] ])?[\[(].*)$`)
	reSeason   = regexp.MustCompile(`^S?([0-9]+)$`)
	reOrdinal  = regexp.MustCompile(`^([0-9]+)(?:ST|ND|RD|TH)$`)
	reRoman    = regexp.MustCompile(`^[IVX]+$`)
	reDigits   = regexp.MustCompile(`(\b|-)[0-9]+(\b|-)`)
	romanLib   = roman.NewRoman()
)

func (r *Remote) IsValid(p string) bool {
	if r.Path == "" {
		return true
	}
	return strings.HasPrefix(p, r.Path)
}

func (r *Remote) FullPath(p string, ignorePath bool) string {
	var remotePath string
	remoteName := r.Name
	if ignorePath {
		remotePath = p
	} else {
		remotePath = path.Join(r.Path, p)
	}
	if remoteName != "" {
		remoteName += ":"
		if remotePath == "." {
			remotePath = ""
		}
	}
	return remoteName + remotePath
}

func SendFileIdToWorker(fileId int64) {
	callbackMu.Lock()
	defer callbackMu.Unlock()

	fileInfo := Put.GetFileInfo(fileId)
	if fileInfo == nil {
		return
	}
	go Put.CleanupTransfers()

	if delayBeforeTransfer > 0 {
		go func() {
			notification.Send(fmt.Sprintf("%s downloaded, transfer will begin in %s", fileInfo.Name, delayBeforeTransfer))
			time.Sleep(delayBeforeTransfer)
			// Get file info from Put.io again in case the file was trashed
			fileInfo = Put.GetFileInfo(fileId)
			if fileInfo == nil {
				return
			}
			taskChan <- fileInfo
		}()
	} else {
		taskChan <- fileInfo
	}
}

func worker() {
	defer workerWg.Done()
	for fileInfo := range taskChan {
		workerWg.Add(1)
		if fileInfo.IsDir {
			go moveFolder(fileInfo)
		} else {
			go moveFile(fileInfo)
		}
	}
}

func checkBeforeTransfer(info *putio.FileInfo) bool {
	if !remoteSrc.IsValid(info.FullPath) {
		notification.Send(fmt.Sprintf("%s downloaded", info.Name))
		return false
	}
	return true
}

func moveFolder(folder *putio.FileInfo) {
	folderMu.Lock()
	defer func() {
		workerWg.Done()
		folderMu.Unlock()
	}()

	if !checkBeforeTransfer(folder) {
		log.Printf("Folder %s skipped", folder.Name)
		return
	}

	if folder.Size > 0 {
		log.Printf("Moving folder %s...", folder.Name)

		src := remoteSrc.FullPath(folder.FullPath, true)
		dest := remoteDest.FullPath(folder.Name, false)
		if rcMoveDir(src, dest) {
			Put.DeleteFile(folder.ID)
			notification.Send(fmt.Sprintf("%s moved", folder.Name))
		} else {
			notification.Send(fmt.Sprintf("error occurred, %s was not moved", folder.Name))
		}
	} else {
		Put.DeleteFile(folder.ID)
	}
}

func moveFile(file *putio.FileInfo) {
	defer workerWg.Done()

	if checkBeforeTransfer(file) {
		log.Printf("Moving file %s...", file.Name)
	} else {
		log.Printf("File %s skipped", file.Name)
		return
	}

	newFilename := file.Name
	if strings.HasPrefix(file.ContentType, putio.ContentTypeVideo) {
		switch renamingStyle {
		case RenamingStyleAnime:
			newFilename = RenameFileInAnimeStyle(file.Name)
		case RenamingStyleTv:
			newFilename = RenameFileInTvStyle(file.Name)
		}
	}

	src := remoteSrc.FullPath(file.FullPath, true)
	dest := remoteDest.FullPath(newFilename, false)
	if !rcMoveFile(src, dest, file.Size) {
		notification.Send(fmt.Sprintf("error occurred, %s was not moved", file.Name))
	} else if file.Name == newFilename {
		notification.Send(fmt.Sprintf("%s moved", file.Name))
	} else {
		notification.Send(fmt.Sprintf("%s moved and renamed", file.Name))
	}
}

func ParseEpisodeInfo(filename string, keepSeason bool) *EpisodeInfo {
	filename = strings.ReplaceAll(filename, "] [", "][")
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
	var lastElem, secondLastElem string
	lastElem = strings.ToUpper(elems[lenElems-1])
	if len(elems) >= 2 {
		secondLastElem = strings.ToUpper(elems[lenElems-2])
	}

	season := 1
	var seasonLength int
	if lastElem == "SEASON" {
		if matches := reOrdinal.FindStringSubmatch(secondLastElem); matches != nil {
			season, _ = strconv.Atoi(matches[1])
			seasonLength = 2
		}
	} else if matches := reSeason.FindStringSubmatch(lastElem); matches != nil {
		season, _ = strconv.Atoi(matches[1])
		if secondLastElem == "PART" || secondLastElem == "SEASON" {
			seasonLength = 2
		} else {
			seasonLength = 1
		}
	} else if matches := reRoman.FindStringSubmatch(lastElem); matches != nil {
		season = romanLib.ToNumber(lastElem)
		seasonLength = 1
	}
	if seasonLength > 0 && !keepSeason {
		elems = elems[:lenElems-seasonLength]
	}
	if matches[3] == "SP" || matches[6] != "" {
		info.Season = 0
	} else if season < 100 {
		info.Season = season
	}

	info.Show = strings.ReplaceAll(strings.Join(elems, " "), "_", " ")
	return info
}

func RenameFileInAnimeStyle(filename string) string {
	info := ParseEpisodeInfo(filename, false)
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
	info := ParseEpisodeInfo(filename, true)
	if info == nil {
		return filename
	}
	extra := info.Extra
	if strings.HasPrefix(extra, "END ") {
		extra = strings.Replace(extra, "END ", "", 1)
	}
	if strings.HasPrefix(extra, "(") {
		extra = " " + extra
	}
	return fmt.Sprintf("%s - S%02dE%s - %s%s", info.Show, info.Season, info.Episode, info.Group, extra)
}
