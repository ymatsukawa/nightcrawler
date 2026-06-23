package match_case

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadingWildcardLike(t *testing.T) {
	type Test struct {
		Input  string
		Expect bool
	}

	tests := []Test{
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="Hello world"`,
			Expect: false,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where name like '%abc';"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where name like %abc;"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where name like 'abc%';"`,
			Expect: false,
		},
	}

	for _, tt := range tests {
		result := HeadingWildcardLike(tt.Input)
		assert.Equal(t, tt.Expect, result)
	}
}

func TestUsingOr(t *testing.T) {
	type Test struct {
		Input  string
		Expect bool
	}

	tests := []Test{
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where a = 1 or b = 2 or c = 3;"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where a = 1 or b = 2;"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where a = 1;"`,
			Expect: false,
		},
	}

	for _, tt := range tests {
		result := UsingOr(tt.Input)
		assert.Equal(t, tt.Expect, result)
	}
}
