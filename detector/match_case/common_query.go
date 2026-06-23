package match_case

import "regexp"

var (
	selectQuery = regexp.MustCompile(`select\s+.*\s+from`)
	insertQuery = regexp.MustCompile(`insert\s+into`)
	updateQuery = regexp.MustCompile(`update\s+.*set`)
	deleteQuery = regexp.MustCompile(`delete\s+from`)
)
