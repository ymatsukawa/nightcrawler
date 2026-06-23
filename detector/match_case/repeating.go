package match_case

import (
	"regexp"
	"strings"
)

var (
	sqlStringValue = regexp.MustCompile(`'[^']*'`)
	sqlNumberValue = regexp.MustCompile(`\b\d+(\.\d+)?\b`)
)

func Repeating(currentLine string, previousLine string) bool {
	if strings.TrimSpace(previousLine) == "" {
		return false
	}

	current := replaceSqlValueToSymbol(currentLine)
	previous := replaceSqlValueToSymbol(previousLine)

	if !isSqlQuery(current) {
		return false
	}

	return current == previous
}

func replaceSqlValueToSymbol(line string) string {
	s := strings.ToLower(line)
	s = sqlStringValue.ReplaceAllLiteralString(s, `'$'`)
	s = sqlNumberValue.ReplaceAllLiteralString(s, `$`)

	return s
}

func isSqlQuery(line string) bool {
	return selectQuery.MatchString(line) ||
		insertQuery.MatchString(line) ||
		updateQuery.MatchString(line) ||
		deleteQuery.MatchString(line)
}
