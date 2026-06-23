package match_case

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsSubquery(t *testing.T) {
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
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where id in (select id from others);"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where id = 1;"`,
			Expect: false,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where price > (  select avg(price) from examples);"`,
			Expect: true,
		},
	}

	for _, tt := range tests {
		result := ContainsSubquery(tt.Input)
		assert.Equal(t, tt.Expect, result)
	}
}
