package detector

import (
	"slices"
	"strings"

	"github.com/ymatsukawa/slow_query/category"
	mc "github.com/ymatsukawa/slow_query/detector/match_case"
)

type ParseInfo struct {
	PreviousLine string
	Suppress     []string
}

type rule struct {
	match    func(string) bool
	category string
}

var detectors = []rule{
	{mc.ContainsSelectAsterisk, category.SelectMany},
	{mc.NoLimit, category.SelectMany},

	{mc.HeadingWildcardLike, category.IndexNoWorth},
	{mc.UsingOr, category.IndexNoWorth},

	{mc.ContainsSubquery, category.UsingSubquery},

	{mc.ManyIn, category.HeavyClause},
	{mc.ManyJoin, category.HeavyClause},

	{mc.SelectCount, category.UsingFunction},
	{mc.WhereDate, category.UsingFunction},
}

func CatchSlowQuery(log string, parseInfo ParseInfo) (string, bool) {
	l := strings.ToLower(log)

	if mc.Repeating(l, parseInfo.PreviousLine) {
		return category.NPlusOne, true
	}

	for _, d := range detectors {
		if slices.Contains(parseInfo.Suppress, d.category) {
			continue
		}
		if d.match(l) {
			return d.category, true
		}
	}

	return "", false
}
