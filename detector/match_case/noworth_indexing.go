package match_case

import "regexp"

var (
	headingWildCardLike = regexp.MustCompile(`like\s+['"]?%`)
	orClause            = regexp.MustCompile(`\bor\b`)
	usingNotClause      = regexp.MustCompile(`\bnot\b`)
	usingIsNull         = regexp.MustCompile(`is\s+null`)
)

func HeadingWildcardLike(log string) bool {
	return headingWildCardLike.MatchString(log)
}

func UsingOr(log string) bool {
	return orClause.MatchString(log)
}

func UsingNot(log string) bool {
	return usingNotClause.MatchString(log)
}

func UsingIsNull(log string) bool {
	return usingIsNull.MatchString(log)
}
