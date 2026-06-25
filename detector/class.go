package detector

import (
	mc "github.com/ymatsukawa/slow_query/detector/match_case"
)

type Class interface {
	GetId() int
	ToSummary() string
	IsMatch(log string) bool
}

type class struct {
	Id      int
	Summary string
	Match   func(log string) bool
}

const (
	SelectMany = iota
	NoLimit
	ManyJoin
	ManyIn
	HeadingWildcardLike
	UsingOr
	UsingNot
	UsingIsNull
	UsingSubquery
	UsingCount
	UsingMax
	UsingMin
	UsingAvg
	UsingSum
	WhereDate
)

const (
	nPlusOne = "N+1"
)

var classes = []Class{
	&class{
		Id:      SelectMany,
		Summary: "SELECT *",
		Match:   mc.ContainsSelectAsterisk,
	},
	&class{
		Id:      NoLimit,
		Summary: "no LIMIT",
		Match:   mc.NoLimit,
	},
	&class{
		Id:      ManyJoin,
		Summary: "many JOIN 2 >=",
		Match:   mc.ManyJoin,
	},
	&class{
		Id:      ManyIn,
		Summary: "many IN values 100 >=",
		Match:   mc.ManyIn,
	},
	&class{
		Id:      HeadingWildcardLike,
		Summary: "LIKE %...",
		Match:   mc.HeadingWildcardLike,
	},
	&class{
		Id:      UsingOr,
		Summary: "using OR",
		Match:   mc.UsingOr,
	},
	&class{
		Id:      UsingNot,
		Summary: "using NOT",
		Match:   mc.UsingNot,
	},
	&class{
		Id:      UsingIsNull,
		Summary: "using IS NULL",
		Match:   mc.UsingIsNull,
	},
	&class{
		Id:      UsingSubquery,
		Summary: "using Subquery",
		Match:   mc.ContainsSubquery,
	},
	&class{
		Id:      UsingCount,
		Summary: "using COUNT",
		Match:   mc.SelectCount,
	},
	&class{
		Id:      UsingMax,
		Summary: "using MAX",
		Match:   mc.SelectMax,
	},
	&class{
		Id:      UsingMin,
		Summary: "using MIN",
		Match:   mc.SelectMin,
	},
	&class{
		Id:      UsingAvg,
		Summary: "using AVG",
		Match:   mc.SelectAvg,
	},
	&class{
		Id:      UsingSum,
		Summary: "using SUM",
		Match:   mc.SelectSum,
	},
	&class{
		Id:      WhereDate,
		Summary: "using WHERE with DATE func",
		Match:   mc.WhereDate,
	},
}

func (c *class) GetId() int {
	return c.Id
}

func (c *class) ToSummary() string {
	return c.Summary
}

func (c *class) IsMatch(log string) bool {
	return c.Match(log)
}
