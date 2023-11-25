package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SQLStringValue string

func ToSQLStringValue(lines ...string) []SQLStringValue {
	var strValues = make([]SQLStringValue, 0, len(lines))
	for _, line := range lines {
		strValues = append(strValues, SQLStringValue(line))
	}
	return strValues
}

type queue []any

func (q *queue) pop() any {
	if len(*q) <= 0 {
		return nil
	}
	val := (*q)[0]

	*q = (*q)[1:]
	return val
}

func SQLReplaceArgs(query string, args ...any) string {
	var argsQueue queue = args

	var resultQuery strings.Builder
	for _, r := range query {
		if r == '?' {
			arg := argsQueue.pop()
			if arg == nil {
				resultQuery.WriteString("?")
				continue
			}
			switch val := arg.(type) {
			case []string:
				resultQuery.WriteString(strings.Join(val, ","))
			case string:
				resultQuery.WriteString(val)
			case []SQLStringValue:
				resultQuery.WriteString(joinQuotedSQLStrings(val, "'", ","))
			case SQLStringValue:
				resultQuery.WriteString("'")
				resultQuery.WriteString(string(val))
				resultQuery.WriteString("'")
			case []int:
				resultQuery.WriteString(JoinQuotedInt(val, "", ","))
			case int:
				resultQuery.WriteString(strconv.Itoa(val))
			case int32:
				resultQuery.WriteString(strconv.FormatInt(int64(val), 10))
			case int64:
				resultQuery.WriteString(strconv.FormatInt(val, 10))
			case uint:
				resultQuery.WriteString(strconv.FormatUint(uint64(val), 10))
			case uint32:
				resultQuery.WriteString(strconv.FormatUint(uint64(val), 10))
			case uint64:
				resultQuery.WriteString(strconv.FormatUint(val, 10))
			case time.Time:
				resultQuery.WriteString("'")
				resultQuery.WriteString(val.Format("2006-01-02 15:04:05"))
				resultQuery.WriteString("'")
			default:
				resultQuery.WriteString(fmt.Sprintf("%v", val))
			}
		} else {
			resultQuery.WriteString(string(r))
		}
	}

	return resultQuery.String()
}

func joinQuotedSQLStrings(lines []SQLStringValue, quote, separator string) string {
	var joined strings.Builder
	for i, line := range lines {
		if i > 0 {
			joined.WriteString(separator)
		}
		joined.WriteString(quote)
		joined.WriteString(string(line))
		joined.WriteString(quote)
	}
	return joined.String()
}
