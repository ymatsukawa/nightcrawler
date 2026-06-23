package match_case

import (
	"regexp"
	"strings"
)

var (
	selectAsterisk = regexp.MustCompile(`select\s+\*`)
)

func ContainsSelectAsterisk(log string) bool {
	return selectAsterisk.MatchString(log)
}

func NoLimit(log string) bool {
	return selectQuery.MatchString(log) && !strings.Contains(log, "limit")
}
