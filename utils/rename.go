package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/chonla/roman-number-go"
)

var (
	reGroup    = regexp.MustCompile(`^\[.+?]`)
	reEpisode  = regexp.MustCompile(`(?i)^\[?(?:EP|#|第)?(SP|OVA|OAD|EX)?([0-9]{2,}(?:\.[0-9])?)?話?(?:v([0-9]))?(\(.+\))?]?$`)
	reSeason   = regexp.MustCompile(`((?:(?i)Season |Part )?([0-9]+|[IVX]+)$)|(([0-9]+)(?:nd|rd|th)?(?:(?i) Season)?$)|(S([0-9]+)$)`)
	reDigits   = regexp.MustCompile(`(\b|-)[0-9]+(\b|-)`)
	reBrackets = regexp.MustCompile(`[\[\]]`)
	romanLib   = roman.NewRoman()
)

func FirstOrElse[T any](predicate func(arg T) bool, defaultVal T, args ...T) T {
	for _, arg := range args {
		if predicate(arg) {
			return arg
		}
	}
	return defaultVal
}

func ParseEpisodeInfo(filename string, keepSeason bool) *EpisodeInfo {
	group := reGroup.FindString(filename)
	if group == "" {
		return nil
	}
	info := &EpisodeInfo{
		Group: group,
	}

	var showParts []string
	var holdParts []string

	name := filename[len(group):]
	if !strings.HasPrefix(name, "[") {
		idx := strings.IndexRune(name, '[')
		if idx > 0 {
			name = name[0:idx] + " " + name[idx:]
		}
	}
	fields := strings.Fields(strings.ReplaceAll(name, "][", "] ["))
	for _, field := range fields {
		if field == "-" {
			continue
		}

		eMatches := reEpisode.FindStringSubmatch(field)
		if eMatches == nil || (eMatches[1] == "" && eMatches[2] == "") {
			if info.Episode == "" {
				showParts = append(showParts, field)
			} else {
				holdParts = append(holdParts, field)
			}
		} else {
			if eMatches[1] != "" || eMatches[4] != "" {
				info.Season = 0
			} else {
				info.Season = 1
			}
			if eMatches[3] != "" {
				info.Version, _ = strconv.Atoi(eMatches[3])
			} else {
				info.Version = 1
			}
			if info.Episode == "" {
				holdParts = append(holdParts, field)
			} else {
				showParts = append(showParts, holdParts...)
				holdParts = []string{field}
			}
			if eMatches[2] != "" {
				info.Episode = eMatches[2]
			} else {
				info.Episode = "01"
			}
		}
	}
	if info.Episode == "" {
		return nil
	}
	if len(showParts) > 0 {
		for i, part := range showParts {
			showParts[i] = strings.ReplaceAll(reBrackets.ReplaceAllString(part, ""), "_", " ")
		}
		info.Show = strings.Join(showParts, " ")
	}
	if info.Show == "" {
		return nil
	}
	for info.Season != 0 {
		sMatches := reSeason.FindStringSubmatch(info.Show)
		if sMatches == nil {
			break
		}
		seasonStr := FirstOrElse[string](func(arg string) bool {
			return arg != ""
		}, "", sMatches[2], sMatches[4], sMatches[6])
		season, err := strconv.Atoi(seasonStr)
		if err != nil {
			season = romanLib.ToNumber(seasonStr)
		}
		if season >= 100 {
			break
		}
		info.Season = season
		if keepSeason {
			break
		}
		idx := strings.LastIndex(info.Show, " "+sMatches[0])
		info.Show = info.Show[:idx]
		break
	}
	if len(holdParts) > 1 {
		for i, part := range holdParts {
			if i <= 1 {
				continue
			}
			if !strings.HasPrefix(part, "[") {
				holdParts[i] = " " + part
			}
		}
		info.Extra = strings.Join(holdParts[1:], "")
	}
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
