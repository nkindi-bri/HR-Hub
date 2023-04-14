package structure

import (
	"fmt"
	"strconv"
)

// FormatLimitOffset returns a SQL string for a given limit & offset.
// Clauses are only added if limit and/or offset are greater than zero.
func FormatLimitOffset(limit, offset uint64) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)
	} else if offset > 0 {
		return fmt.Sprintf(`OFFSET %d`, offset)
	}
	return ""
}

func FormatPeriod(from, to string) string {
	if from == "" && to == "" {
		return ""
	}
	return fmt.Sprintf("AND BETWEEN %s AND %s", from, to)
}

// PosArg returns a SQL positional argument and increments the position.
func PosArg(pos *int) string {
	(*pos)++
	return fmt.Sprintf("$%d", *pos)
}

//ConvertToInt64 converts a string to an int64
func ConvertToInt(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

//ConvertToFloat64 converts a string to an float64
func ConvertToFloat(str string) float64 {
	i, _ := strconv.ParseFloat(str, 64)
	return i
}
