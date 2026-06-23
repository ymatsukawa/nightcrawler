package match_case

import "regexp"

var (
	headingWildCardLike = regexp.MustCompile(`like\s+['"]?%`)
	orClause            = regexp.MustCompile(`\bor\b`)
)

func HeadingWildcardLike(log string) bool {
	return headingWildCardLike.MatchString(log)
}

func UsingOr(log string) bool {
	return orClause.MatchString(log)
}
