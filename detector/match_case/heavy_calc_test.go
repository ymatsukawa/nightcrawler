package match_case

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManyIn(t *testing.T) {
	type Test struct {
		Input  string
		Expect bool
	}

	manyValues := strings.Repeat("1,", manyInThreshold-1) + "1"

	tests := []Test{
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where id in (` + manyValues + `);"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where id in (1, 2, 3);"`,
			Expect: false,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where id in ();"`,
			Expect: false,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples;"`,
			Expect: false,
		},
	}

	for _, tt := range tests {
		result := ManyIn(tt.Input)
		assert.Equal(t, tt.Expect, result)
	}
}

func TestManyJoin(t *testing.T) {
	type Test struct {
		Input  string
		Expect bool
	}

	tests := []Test{
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from a join b on a.id = b.a_id join c on c.b_id = b.id;"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from a join b on a.id = b.a_id;"`,
			Expect: false,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples;"`,
			Expect: false,
		},
	}

	for _, tt := range tests {
		result := ManyJoin(tt.Input)
		assert.Equal(t, tt.Expect, result)
	}
}
