package nightcrawler

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	d "github.com/ymatsukawa/nightcrawler/detector"
)

func TestNewSlogHandler(t *testing.T) {
	type Args struct {
		SourceHandler slog.Handler
		Suppressors   []d.Suppressor
	}

	tests := []struct {
		Name string
		Args Args
	}{
		{
			Name: "Create new slog handler with suppress",
			Args: Args{
				SourceHandler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
				Suppressors:   []d.Suppressor{},
			},
		},
		{
			Name: "Create new slog handler with suppress (with categories)",
			Args: Args{
				SourceHandler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
				Suppressors:   []d.Suppressor{d.SelectMany},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			handler := NewSlogHandler(tt.Args.SourceHandler, tt.Args.Suppressors)
			assert.NotNil(t, handler)
			assert.Equal(t, tt.Args.SourceHandler, handler.SourceHandler)
			assert.Equal(t, tt.Args.Suppressors, handler.LogInfo.Suppressors)
		})
	}
}

func TestLogging(t *testing.T) {
	type TestCase struct {
		Name                     string
		Log                      string
		ExpectSlowQueryAttrShown bool
		ExpectAttrValue          string
	}

	manyInLog := "select id from users where id in (" + strings.Repeat("1, ", 100) + "1) limit 10"

	tests := []TestCase{
		{
			Name:                     "No slow query",
			Log:                      "select id from users limit 10",
			ExpectSlowQueryAttrShown: false,
			ExpectAttrValue:          "",
		},
		{
			Name:                     "ContainsSelectAsterisk: select *",
			Log:                      "select * from users limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "SELECT *",
		},
		{
			Name:                     "NoLimit: select without limit",
			Log:                      "select id from users",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "no LIMIT",
		},
		{
			Name:                     "ManyJoin: two or more joins",
			Log:                      "select u.id from users u join orders o on o.uid = u.id join items i on i.oid = o.id limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "many JOIN 2 >=",
		},
		{
			Name:                     "ManyIn: in clause with 100+ values",
			Log:                      manyInLog,
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "many IN values 100 >=",
		},
		{
			Name:                     "HeadingWildcardLike: leading wildcard like",
			Log:                      "select id from users where name like '%john%' limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "LIKE %...",
		},
		{
			Name:                     "UsingOr: or clause",
			Log:                      "select id from users where age = 1 or age = 2 limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "using OR",
		},
		{
			Name:                     "UsingNot: not clause",
			Log:                      "select id from users where status not in (1, 2, 3) limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "using NOT",
		},
		{
			Name:                     "UsingIsNull: is null clause",
			Log:                      "select id from users where deleted_at is null limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "using IS NULL",
		},
		{
			Name:                     "ContainsSubquery: subquery in where",
			Log:                      "select id from users where id in (select uid from orders) limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "using Subquery",
		},
		{
			Name:                     "SelectCount: count function",
			Log:                      "select count(id) from products limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "using COUNT",
		},
		{
			Name:                     "SelectMax: max function",
			Log:                      "select max(price) from products limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "using MAX",
		},
		{
			Name:                     "SelectMin: min function",
			Log:                      "select min(price) from products limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "using MIN",
		},
		{
			Name:                     "SelectAvg: avg function",
			Log:                      "select avg(price) from products limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "using AVG",
		},
		{
			Name:                     "SelectSum: sum function",
			Log:                      "select sum(price) from products limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "using SUM",
		},
		{
			Name:                     "WhereDate: date function in where",
			Log:                      "select id from products where date(created_at) = '2020-01-01' limit 10",
			ExpectSlowQueryAttrShown: true,
			ExpectAttrValue:          "using WHERE with DATE func",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var buf bytes.Buffer
			baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
			logger := slog.New(NewSlogHandler(baseHandler, nil))

			logger.Info(tt.Log)

			logOutput := buf.String()
			if tt.ExpectSlowQueryAttrShown {
				assert.Contains(t, logOutput, `slow_query="`+tt.ExpectAttrValue+`"`)
			} else {
				assert.NotContains(t, logOutput, `slow_query="`)
			}
		})
	}
}

func TestLoggingRepeating(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(NewSlogHandler(baseHandler, nil))

	logger.Info("select name from members where id = 1")
	buf.Reset()
	logger.Info("select name from members where id = 2")

	assert.Contains(t, buf.String(), "slow_query=N+1")
}

func TestLoggingWithSuppress(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(NewSlogHandler(baseHandler, []d.Suppressor{d.UsingCount}))

	logger.Info("select COUNT(*) from users limit 100")

	assert.NotContains(t, buf.String(), "slow_query=")
}
