package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
