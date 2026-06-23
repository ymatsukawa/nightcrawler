package match_case

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectCount(t *testing.T) {
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
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select count(*) from examples;"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select count (id) from examples;"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select id from examples;"`,
			Expect: false,
		},
	}

	for _, tt := range tests {
		result := SelectCount(tt.Input)
		assert.Equal(t, tt.Expect, result)
	}
}

func TestWhereDate(t *testing.T) {
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
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where date(created_at) = '2026-06-22';"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select date(created_at) from examples;"`,
			Expect: false,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where id = 1;"`,
			Expect: false,
		},
	}

	for _, tt := range tests {
		result := WhereDate(tt.Input)
		assert.Equal(t, tt.Expect, result)
	}
}
