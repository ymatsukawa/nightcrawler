package match_case

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsSelectAsterisk(t *testing.T) {
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
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples;"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples;"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * \n\tfrom examples;"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select *     from examples;"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select id, name, email from examples;"`,
			Expect: false,
		},
	}

	for _, tt := range tests {
		result := ContainsSelectAsterisk(tt.Input)
		assert.Equal(t, tt.Expect, result)
	}
}

func TestNoLimit(t *testing.T) {
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
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples;"`,
			Expect: true,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples limit 10;"`,
			Expect: false,
		},
		{
			Input:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="insert into examples (id) values (1);"`,
			Expect: false,
		},
	}

	for _, tt := range tests {
		result := NoLimit(tt.Input)
		assert.Equal(t, tt.Expect, result)
	}
}
