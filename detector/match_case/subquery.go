package match_case

import "regexp"

var whereSubquery = regexp.MustCompile(`where\b[\s\S]*\(\s*select\b`)

func ContainsSubquery(log string) bool {
	return whereSubquery.MatchString(log)
}
