package utils

import (
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func JoinQuoted(lines []string, quote, separator string) string {
	var joined strings.Builder
	for i, line := range lines {
		if i > 0 {
			joined.WriteString(separator)
		}
		joined.WriteString(quote)
		joined.WriteString(line)
		joined.WriteString(quote)
	}
	return joined.String()
}

func JoinQuotedInt(lines []int, quote, separator string) string {
	var joined strings.Builder
	for i, line := range lines {
		if i > 0 {
			joined.WriteString(separator)
		}
		joined.WriteString(quote)
		joined.WriteString(strconv.Itoa(line))
		joined.WriteString(quote)
	}
	return joined.String()
}

func FuncName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	funcNameBeginning := strings.LastIndex(frame.Function, ".") + 1
	if 0 < funcNameBeginning && funcNameBeginning < len(frame.Function) {
		return ToUnderscoreACII(frame.Function[funcNameBeginning:])
	}
	return "unknown_function"
}

var (
	regexpUpperWord        = regexp.MustCompile(`[A-Z][a-z]+`)
	regexpAbbreviation     = regexp.MustCompile(`[A-Z]+`) // must be applied after replacing all simple words
	regexpSpace            = regexp.MustCompile(`[\n\r ]`)
	regexpMultiUnderscores = regexp.MustCompile(`_{2,}`)
)

func ToUnderscoreACII(s string) (underscored string) {
	underscored = regexpUpperWord.ReplaceAllStringFunc(s, func(upperWord string) string {
		return "_" + strings.ToLower(upperWord)
	})
	underscored = regexpAbbreviation.ReplaceAllStringFunc(underscored, func(abbreviation string) string {
		return "_" + strings.ToLower(abbreviation)
	})
	underscored = regexpSpace.ReplaceAllString(underscored, "_")
	underscored = regexpMultiUnderscores.ReplaceAllString(underscored, "_")
	underscored = strings.Trim(underscored, "_")
	return strings.ToLower(underscored)
}

func JoinParams(first, second string) string {
	return first + "#" + second
}

func GetFirst(str string) string {
	idx := strings.Index(str, "#")
	if idx == -1 {
		return str
	}
	return str[:idx]
}

func GetSecond(str string) string {
	idx := strings.Index(str, "#")
	if idx == -1 {
		return ""
	}
	return str[idx+1:]
}

func EscapeChars(str string, chars map[rune]struct{}) string {
	var escapedStr strings.Builder
	for _, r := range str {
		if _, ok := chars[r]; ok {
			escapedStr.WriteString(`\`)
		}
		escapedStr.WriteRune(r)
	}
	return escapedStr.String()
}
