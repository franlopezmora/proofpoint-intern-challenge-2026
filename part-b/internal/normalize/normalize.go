package normalize

import (
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	UnknownAirDate   = "Unknown"
	UntitledEpisode  = "Untitled Episode"
	defaultSeriesCap = ""
)

var supportedDateLayouts = []string{
	"2006-01-02",
	"2006/01/02",
	"02-01-2006",
	"02/01/2006",
	"2006.01.02",
	"Jan 2 2006",
	"January 2 2006",
}

func NormalizeSeriesName(raw string) (string, bool) {
	cleaned, changed := normalizeDisplayText(raw)
	return cleaned, changed
}

func NormalizeEpisodeTitle(raw string) (string, bool) {
	cleaned, changed := normalizeDisplayText(raw)
	if cleaned == "" {
		return UntitledEpisode, true
	}
	return cleaned, changed
}

func NormalizeComparisonText(raw string) string {
	return strings.ToLower(collapseWhitespace(stripControlChars(raw)))
}

func ParseSeasonNumber(raw string) (int, bool) {
	return parseNonNegativeInt(raw)
}

func ParseEpisodeNumber(raw string) (int, bool) {
	return parseNonNegativeInt(raw)
}

func ParseAirDate(raw string) (string, bool) {
	candidate := collapseWhitespace(stripControlChars(raw))
	if candidate == "" {
		return UnknownAirDate, true
	}

	for _, layout := range supportedDateLayouts {
		parsed, err := time.Parse(layout, candidate)
		if err != nil {
			continue
		}
		return parsed.Format("2006-01-02"), parsed.Format("2006-01-02") != candidate
	}

	return UnknownAirDate, true
}

func IsTitleKnown(title string) bool {
	return title != UntitledEpisode
}

func IsAirDateKnown(airDate string) bool {
	return airDate != UnknownAirDate
}

func normalizeDisplayText(raw string) (string, bool) {
	sanitized := collapseWhitespace(stripControlChars(raw))
	if sanitized == "" {
		return defaultSeriesCap, raw != ""
	}

	titleCased := toTitleLike(sanitized)
	changed := titleCased != raw
	return titleCased, changed
}

func parseNonNegativeInt(raw string) (int, bool) {
	candidate := collapseWhitespace(stripControlChars(raw))
	if candidate == "" {
		return 0, true
	}

	number, err := strconv.Atoi(candidate)
	if err != nil || number < 0 {
		return 0, true
	}
	if strconv.Itoa(number) != candidate {
		return 0, true
	}
	return number, number == 0 && candidate != "0"
}

func stripControlChars(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, input)
}

func collapseWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

func toTitleLike(input string) string {
	lowered := strings.ToLower(input)
	tokens := strings.Split(lowered, " ")
	for i, token := range tokens {
		tokens[i] = capitalizeToken(token)
	}
	return strings.Join(tokens, " ")
}

func capitalizeToken(token string) string {
	if token == "" {
		return token
	}

	runes := []rune(token)
	for i, r := range runes {
		if unicode.IsLetter(r) {
			runes[i] = unicode.ToUpper(r)
			break
		}
	}
	return string(runes)
}
