package match_case

import (
	"regexp"
	"strings"
)

const (
	manyInThreshold   = 100
	manyJoinThreshold = 2
)

var (
	inClause  = regexp.MustCompile(`\bin\s*\(([^)]*)\)`)
	joinToken = regexp.MustCompile(`\bjoin\b`)
)

func ManyIn(log string) bool {
	for _, m := range inClause.FindAllStringSubmatch(log, -1) {
		if strings.TrimSpace(m[1]) == "" {
			continue
		}
		if strings.Count(m[1], ",")+1 >= manyInThreshold {
			return true
		}
	}
	return false
}

func ManyJoin(log string) bool {
	return len(joinToken.FindAllString(log, -1)) >= manyJoinThreshold
}
