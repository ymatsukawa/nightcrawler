package match_case

import "regexp"

var (
	countFunction = regexp.MustCompile(`\bcount\s*\(`)
	maxFunction   = regexp.MustCompile(`\bmax\s*\(`)
	minFunction   = regexp.MustCompile(`\bmin\s*\(`)
	avgFunction   = regexp.MustCompile(`\bavg\s*\(`)
	sumFunction   = regexp.MustCompile(`\bsum\s*\(`)
	whereDate     = regexp.MustCompile(`where\b[\s\S]*\bdate\s*\(`)
)

func SelectCount(log string) bool {
	return countFunction.MatchString(log)
}

func SelectMax(log string) bool {
	return maxFunction.MatchString(log)
}

func SelectMin(log string) bool {
	return minFunction.MatchString(log)
}

func SelectAvg(log string) bool {
	return avgFunction.MatchString(log)
}

func SelectSum(log string) bool {
	return sumFunction.MatchString(log)
}

func WhereDate(log string) bool {
	return whereDate.MatchString(log)
}
