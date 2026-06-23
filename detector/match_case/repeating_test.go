package match_case

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeating(t *testing.T) {
	type Test struct {
		Current  string
		Previous string
		Expect   bool
	}

	tests := []Test{
		{
			Current:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where id = 1;"`,
			Previous: `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where id = 2;"`,
			Expect:   true,
		},
		{
			Current:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where name = 'alice';"`,
			Previous: `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where name = 'bob';"`,
			Expect:   true,
		},
		{
			Current:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where id = 1;"`,
			Previous: ``,
			Expect:   false,
		},
		{
			Current:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="Hello world"`,
			Previous: `time=1970-01-01T11:12:00.000+09:00 level=info msg="Hello world"`,
			Expect:   false,
		},
		{
			Current:  `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from examples where id = 1;"`,
			Previous: `time=1970-01-01T11:12:00.000+09:00 level=info msg="select * from others where id = 2;"`,
			Expect:   false,
		},
	}

	for _, tt := range tests {
		result := Repeating(tt.Current, tt.Previous)
		assert.Equal(t, tt.Expect, result)
	}
}
