package detector

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatchSlowQuery(t *testing.T) {
	type Test struct {
		Name          string
		Input         string
		Previous      string
		ExpectSummary string
		ExpectOk      bool
	}

	tests := []Test{
		{
			Name:          "N+1 (repeating)",
			Input:         `select name from users where code = 2`,
			Previous:      `select name from users where code = 1`,
			ExpectSummary: nPlusOne,
			ExpectOk:      true,
		},
		{
			Name:          "Select many (select *)",
			Input:         `select * from users`,
			ExpectSummary: testGetSummary(SelectMany),
			ExpectOk:      true,
		},
		{
			Name:          "Select many (no limit)",
			Input:         `select "users"."id", "users"."name" from users`,
			ExpectSummary: testGetSummary(NoLimit),
			ExpectOk:      true,
		},
		{
			Name:          "Index no worth (leading wildcard like)",
			Input:         `update users set name = 'x' where name like '%test'`,
			ExpectSummary: testGetSummary(HeadingWildcardLike),
			ExpectOk:      true,
		},
		{
			Name:          "Index no worth (using or)",
			Input:         `update users set name = 'x' where a = 1 or b = 2`,
			ExpectSummary: testGetSummary(UsingOr),
			ExpectOk:      true,
		},
		{
			Name:          "Index no worth (using not)",
			Input:         `update users set name = 'x' where a = 1 and not b = 2`,
			ExpectSummary: testGetSummary(UsingNot),
			ExpectOk:      true,
		},
		{
			Name:          "Index no worth (is null)",
			Input:         `update users set name = 'x' where deleted_at is null`,
			ExpectSummary: testGetSummary(UsingIsNull),
			ExpectOk:      true,
		},
		{
			Name:          "Using subquery",
			Input:         `select id from users where id in (select user_id from orders) limit 10`,
			ExpectSummary: testGetSummary(UsingSubquery),
			ExpectOk:      true,
		},
		{
			Name:          "Heavy calc (many join)",
			Input:         `select id from a join b on a.id = b.a_id join c on b.id = c.b_id limit 5`,
			ExpectSummary: testGetSummary(ManyJoin),
			ExpectOk:      true,
		},
		{
			Name:          "Heavy calc (many in)",
			Input:         `select id from users where id in (` + strings.Repeat("1,", 99) + `1) limit 5`,
			ExpectSummary: testGetSummary(ManyIn),
			ExpectOk:      true,
		},
		{
			Name:          "Using function (count)",
			Input:         `select count(*) from users limit 5`,
			ExpectSummary: testGetSummary(UsingCount),
			ExpectOk:      true,
		},
		{
			Name:          "no slow query",
			Input:         `select id from users limit 10`,
			ExpectSummary: "",
			ExpectOk:      false,
		},
		{
			Name:          "Using function (avg)",
			Input:         `select avg(age) from users limit 5`,
			ExpectSummary: testGetSummary(UsingAvg),
			ExpectOk:      true,
		},
		{
			Name:          "Using function (max)",
			Input:         `select max(age) from users limit 5`,
			ExpectSummary: testGetSummary(UsingMax),
			ExpectOk:      true,
		},
		{
			Name:          "Using function (min)",
			Input:         `select min(age) from users limit 5`,
			ExpectSummary: testGetSummary(UsingMin),
			ExpectOk:      true,
		},
		{
			Name:          "Using function (sum)",
			Input:         `select sum(age) from users limit 5`,
			ExpectSummary: testGetSummary(UsingSum),
			ExpectOk:      true,
		},
		{
			Name:          "Using function (where date)",
			Input:         `select id from users where date(created_at) = '2023-01-01' limit 5`,
			ExpectSummary: testGetSummary(WhereDate),
			ExpectOk:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			category, ok := CatchSlowQuery(tt.Input, ParseInfo{PreviousLine: tt.Previous})
			assert.Equal(t, tt.ExpectOk, ok)
			assert.Equal(t, tt.ExpectSummary, category)
		})
	}
}

func TestCatchSlowQueryBySuppressing(t *testing.T) {
	type Test struct {
		Name        string
		Suppress    []int
		Log         string
		ExpectCatch bool
	}

	tests := []Test{
		{
			Name:        "Suppressed select and no limit",
			Suppress:    []int{SelectMany, NoLimit},
			Log:         `select * from users limit 1`,
			ExpectCatch: false,
		},
		{
			Name:        "Not Suppressed select many",
			Suppress:    []int{},
			Log:         `select * from users limit 10`,
			ExpectCatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, ok := CatchSlowQuery(tt.Log, ParseInfo{Suppress: tt.Suppress})
			assert.Equal(t, tt.ExpectCatch, ok)
		})
	}
}

func testGetSummary(id int) string {
	for _, c := range classes {
		if c.GetId() == id {
			return c.ToSummary()
		}
	}
	return ""
}
