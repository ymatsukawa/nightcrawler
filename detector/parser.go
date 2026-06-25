package detector

import (
	"slices"
	"strings"

	mc "github.com/ymatsukawa/nightcrawler/detector/match_case"
)

type ParseInfo struct {
	PreviousLine string
	Suppress     []int
}

func CatchSlowQuery(log string, parseInfo ParseInfo) (string, bool) {
	l := strings.ToLower(log)

	if mc.Repeating(l, parseInfo.PreviousLine) {
		return nPlusOne, true
	}

	for _, c := range classes {
		if slices.Contains(parseInfo.Suppress, c.GetId()) {
			continue
		}
		if c.IsMatch(l) {
			return c.ToSummary(), true
		}
	}

	return "", false
}
