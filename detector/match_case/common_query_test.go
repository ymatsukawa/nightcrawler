package match_case

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommonQuery(t *testing.T) {
	type Test struct {
		Input  string
		Query  *regexp.Regexp
		Expect bool
	}

	tests := []Test{
		{
			Input:  `select * from`,
			Query:  selectQuery,
			Expect: true,
		},
		{
			Input:  `select   *    from`,
			Query:  selectQuery,
			Expect: true,
		},
		{
			Input:  `select* from`,
			Query:  selectQuery,
			Expect: false,
		},
		{
			Input:  `select *from`,
			Query:  selectQuery,
			Expect: false,
		},
		{
			Input:  `insert into`,
			Query:  insertQuery,
			Expect: true,
		},
		{
			Input:  `insert   into`,
			Query:  insertQuery,
			Expect: true,
		},
		{
			Input:  `insertinto`,
			Query:  insertQuery,
			Expect: false,
		},
		{
			Input:  `inset into`,
			Query:  insertQuery,
			Expect: false,
		},
		{
			Input:  `update examples set`,
			Query:  updateQuery,
			Expect: true,
		},
		{
			Input:  `update   examples    set`,
			Query:  updateQuery,
			Expect: true,
		},
		{
			Input:  `updateset`,
			Query:  updateQuery,
			Expect: false,
		},
		{
			Input:  `updat set`,
			Query:  updateQuery,
			Expect: false,
		},
		{
			Input:  `delete from examples`,
			Query:  deleteQuery,
			Expect: true,
		},
		{
			Input:  `delete    from examples`,
			Query:  deleteQuery,
			Expect: true,
		},
		{
			Input:  `deletefrom examples`,
			Query:  deleteQuery,
			Expect: false,
		},
		{
			Input:  `delet from examples`,
			Query:  deleteQuery,
			Expect: false,
		},
	}

	for _, tt := range tests {
		result := tt.Query.MatchString(tt.Input)
		assert.Equal(t, result, tt.Expect, tt.Input)
	}
}
