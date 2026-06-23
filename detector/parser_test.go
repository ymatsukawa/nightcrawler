package detector

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ymatsukawa/slow_query/category"
)

func TestCatchSlowQuery(t *testing.T) {
	type Test struct {
		Name           string
		Input          string
		Previous       string
		ExpectCategory string
		ExpectOk       bool
	}

	tests := []Test{
		{
			Name:           "N+1 (repeating)",
			Input:          `select name from users where id = 2`,
			Previous:       `select name from users where id = 1`,
			ExpectCategory: category.NPlusOne,
			ExpectOk:       true,
		},
		{
			Name:           "Select many (select *)",
			Input:          `select * from users`,
			ExpectCategory: category.SelectMany,
			ExpectOk:       true,
		},
		{
			Name:           "Select many (no limit)",
			Input:          `select "users"."id", "users"."name" from users`,
			ExpectCategory: category.SelectMany,
			ExpectOk:       true,
		},
		{
			Name:           "Index no worth (leading wildcard like)",
			Input:          `update users set name = 'x' where name like '%test'`,
			ExpectCategory: category.IndexNoWorth,
			ExpectOk:       true,
		},
		{
			Name:           "Index no worth (using or)",
			Input:          `update users set name = 'x' where a = 1 or b = 2`,
			ExpectCategory: category.IndexNoWorth,
			ExpectOk:       true,
		},
		{
			Name:           "Index no worth (using not)",
			Input:          `update users set name = 'x' where a = 1 and not b = 2`,
			ExpectCategory: category.IndexNoWorth,
			ExpectOk:       true,
		},
		{
			Name:           "Index no worth (is null)",
			Input:          `update users set name = 'x' where deleted_at is null`,
			ExpectCategory: category.IndexNoWorth,
			ExpectOk:       true,
		},
		{
			Name:           "Using subquery",
			Input:          `select id from users where id in (select user_id from orders) limit 10`,
			ExpectCategory: category.UsingSubquery,
			ExpectOk:       true,
		},
		{
			Name:           "Heavy calc (many join)",
			Input:          `select id from a join b on a.id = b.a_id join c on b.id = c.b_id limit 5`,
			ExpectCategory: category.HeavyCalc,
			ExpectOk:       true,
		},
		{
			Name:           "Heavy calc (many in)",
			Input:          `select id from users where id in (` + strings.Repeat("1,", 99) + `1) limit 5`,
			ExpectCategory: category.HeavyCalc,
			ExpectOk:       true,
		},
		{
			Name:           "Using function (count)",
			Input:          `select count(*) from users limit 5`,
			ExpectCategory: category.UsingFunction,
			ExpectOk:       true,
		},
		{
			Name:           "no slow query",
			Input:          `select id from users limit 10`,
			ExpectCategory: "",
			ExpectOk:       false,
		},
		{
			Name:           "Using function (avg)",
			Input:          `select avg(age) from users limit 5`,
			ExpectCategory: category.UsingFunction,
			ExpectOk:       true,
		},
		{
			Name:           "Using function (max)",
			Input:          `select max(age) from users limit 5`,
			ExpectCategory: category.UsingFunction,
			ExpectOk:       true,
		},
		{
			Name:           "Using function (min)",
			Input:          `select min(age) from users limit 5`,
			ExpectCategory: category.UsingFunction,
			ExpectOk:       true,
		},
		{
			Name:           "Using function (sum)",
			Input:          `select sum(age) from users limit 5`,
			ExpectCategory: category.UsingFunction,
			ExpectOk:       true,
		},
		{
			Name:           "Using function (where date)",
			Input:          `select id from users where date(created_at) = '2023-01-01' limit 5`,
			ExpectCategory: category.UsingFunction,
			ExpectOk:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			category, ok := CatchSlowQuery(tt.Input, ParseInfo{PreviousLine: tt.Previous})
			assert.Equal(t, tt.ExpectOk, ok)
			assert.Equal(t, tt.ExpectCategory, category)
		})
	}
}

func TestCatchSlowQueryBySuppressing(t *testing.T) {
	type Test struct {
		Name     string
		Suppress []string
		Log      string
		Expect   bool
	}

	tests := []Test{
		{
			Name:     "Suppressed select many category",
			Suppress: []string{category.SelectMany},
			Log:      `select * from users`,
			Expect:   false,
		},
		{
			Name:     "Not Suppressed select many category",
			Suppress: []string{},
			Log:      `select * from users`,
			Expect:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, ok := CatchSlowQuery(tt.Log, ParseInfo{Suppress: tt.Suppress})
			assert.Equal(t, tt.Expect, ok)
		})
	}
}
