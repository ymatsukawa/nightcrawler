package match_case

import "regexp"

var (
	countFunction = regexp.MustCompile(`\bcount\s*\(`)
	whereDate     = regexp.MustCompile(`where\b[\s\S]*\bdate\s*\(`)
)

func SelectCount(log string) bool {
	return countFunction.MatchString(log)
}

func WhereDate(log string) bool {
	return whereDate.MatchString(log)
}
