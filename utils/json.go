package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Diff struct {
	Path  string
	Left  string
	Right string
	Error error
}

type Diffs []Diff

var (
	ErrLengthsDontMatch = errors.New("lengths don't match")
	ErrTypesDontMatch   = errors.New("types don't match")
	ErrKeysDontMatch    = errors.New("keys don't match")
	ErrValuesDontMatch  = errors.New("values don't match")
)

func CompareJSON(leftJSON, rightJSON string) (diff Diffs, err error) {
	var (
		leftContent  any
		rightContent any
	)
	if err := json.Unmarshal([]byte(leftJSON), &leftContent); err != nil {
		return nil, fmt.Errorf("error unmarshaling first json file: %w", err)
	}
	if err := json.Unmarshal([]byte(rightJSON), &rightContent); err != nil {
		return nil, fmt.Errorf("error unmarshaling second json file: %w", err)
	}

	diff = compare(leftContent, rightContent)

	return diff, nil
}

func compare(content1, content2 any) []Diff {
	var diff []Diff

	// no difference if null
	if content1 == nil && content2 == nil {
		return diff
	}

	if content1 == nil || content2 == nil {
		jsonContent1, err1 := json.Marshal(content1)
		jsonContent2, err2 := json.Marshal(content2)
		jsonContent1Str := string(jsonContent1)
		jsonContent2Str := string(jsonContent2)
		if err1 != nil || err2 != nil {
			jsonContent1Str = fmt.Sprintf("%v", content1)
			jsonContent2Str = fmt.Sprintf("%v", content2)
		}
		return append(diff, Diff{
			Left:  jsonContent1Str,
			Right: jsonContent2Str,
			Error: ErrValuesDontMatch,
		})
	}

	switch val1 := content1.(type) {
	case map[string]any:
		val2, ok := content2.(map[string]any)
		if !ok {
			return append(diff, Diff{
				Left:  fmt.Sprintf("%T", content1),
				Right: fmt.Sprintf("%T", content2),
				Error: ErrTypesDontMatch,
			})
		}
		if len(val1) != len(val2) {
			diff = append(diff, Diff{
				Left:  strconv.Itoa(len(val1)),
				Right: strconv.Itoa(len(val2)),
				Error: ErrLengthsDontMatch,
			})
		}
		for _, k := range uniqueKeys(val1, val2) {
			if v1, ok := val1[k]; !ok {
				diff = append(diff, Diff{
					Right: k,
					Error: ErrKeysDontMatch,
				})
			} else if v2, ok := val2[k]; !ok {
				diff = append(diff, Diff{
					Left:  k,
					Error: ErrKeysDontMatch,
				})
			} else if vDiff := compare(v1, v2); len(vDiff) > 0 {
				for i := range vDiff {
					if vDiff[i].Error != nil {
						vDiff[i].Path = addPath(vDiff[i].Path, k)
					}
				}
				diff = append(diff, vDiff...)
			}
		}
		return diff

	case []any:
		val2, ok := content2.([]any)
		if !ok {
			return append(diff, Diff{
				Left:  fmt.Sprintf("%T", content1),
				Right: fmt.Sprintf("%T", content2),
				Error: ErrTypesDontMatch,
			})
		}
		if len(val1) != len(val2) {
			diff = append(diff, Diff{
				Left:  strconv.Itoa(len(val1)),
				Right: strconv.Itoa(len(val2)),
				Error: ErrLengthsDontMatch,
			})
		}
		for i := 0; i < longestLen(val1, val2); i++ {
			if vDiff := compare(tryGetElem(val1, i), tryGetElem(val2, i)); len(vDiff) > 0 {
				for j := range vDiff {
					if vDiff[j].Error != nil {
						vDiff[j].Path = addPath(vDiff[j].Path, strconv.Itoa(i))
					}
				}
				diff = append(diff, vDiff...)
			}
		}
		return diff

	case float64:
		if val2, ok := content2.(float64); !ok {
			return append(diff, Diff{
				Left:  fmt.Sprintf("%T", content1),
				Right: fmt.Sprintf("%T", content2),
				Error: ErrTypesDontMatch,
			})
		} else if val1 != val2 {
			diff = append(diff, Diff{
				Left:  strconv.FormatFloat(val1, 'f', -1, 64),
				Right: strconv.FormatFloat(val2, 'f', -1, 64),
				Error: ErrValuesDontMatch,
			})
		}
		return diff

	case bool:
		if val2, ok := content2.(bool); !ok {
			return append(diff, Diff{
				Left:  fmt.Sprintf("%T", content1),
				Right: fmt.Sprintf("%T", content2),
				Error: ErrTypesDontMatch,
			})
		} else if val1 != val2 {
			diff = append(diff, Diff{
				Left:  strconv.FormatBool(val1),
				Right: strconv.FormatBool(val2),
				Error: ErrValuesDontMatch,
			})
		}
		return diff

	case string:
		if val2, ok := content2.(string); !ok {
			return append(diff, Diff{
				Left:  fmt.Sprintf("%T", content1),
				Right: fmt.Sprintf("%T", content2),
				Error: ErrTypesDontMatch,
			})
		} else if val1 != val2 {
			diff = append(diff, Diff{
				Left:  val1,
				Right: val2,
				Error: ErrValuesDontMatch,
			})
		}
		return diff
	}

	return diff
}

func longestLen(c1, c2 []any) int {
	if len(c1) > len(c2) {
		return len(c1)
	}
	return len(c2)
}

func tryGetElem(arr []any, idx int) any {
	if idx >= len(arr) || idx < 0 {
		return nil
	} else {
		return arr[idx]
	}
}

func addPath(current string, add string) string {
	var separator string
	if current != "" {
		separator = "."
	}
	return add + separator + current
}

func uniqueKeys(left, right map[string]any) []string {
	keys := NewSet[string]()
	for k := range left {
		keys.Add(k)
	}
	for k := range right {
		keys.Add(k)
	}
	return keys.Values()
}

func (ds Diffs) String() string {
	var prettyPrint strings.Builder

	pathDiffs := groupByPath(ds)

	for path, pDiffs := range pathDiffs {
		var writeErr = true
		prettyPrint.WriteString("#" + path + ":\n")
		for _, diff := range pDiffs {
			if diff.Error == ErrLengthsDontMatch {
				if writeErr {
					prettyPrint.WriteString(ErrLengthsDontMatch.Error())
					writeErr = false
				}
				left, _ := strconv.Atoi(diff.Left)
				right, _ := strconv.Atoi(diff.Right)
				lenDiff := left - right
				var lenDiffStr string
				if lenDiff >= 0 {
					lenDiffStr = "+" + strconv.Itoa(lenDiff)
				} else {
					lenDiffStr = strconv.Itoa(lenDiff)
				}
				prettyPrint.WriteString(" (" + lenDiffStr + ")\n")
			}
		}

		writeErr = true
		for _, diff := range pDiffs {
			if diff.Error == ErrTypesDontMatch {
				if writeErr {
					prettyPrint.WriteString(ErrTypesDontMatch.Error() + "\n")
					writeErr = false
				}
				prettyPrint.WriteString("+" + diff.Left + "\n")
				prettyPrint.WriteString("-" + diff.Right + "\n")
			}
		}

		writeErr = true
		for _, diff := range pDiffs {
			if diff.Error == ErrKeysDontMatch {
				if writeErr {
					prettyPrint.WriteString(ErrKeysDontMatch.Error() + "\n")
					writeErr = false
				}
				if diff.Left != "" {
					prettyPrint.WriteString("+" + diff.Left + "\n")
				}
				if diff.Right != "" {
					prettyPrint.WriteString("-" + diff.Right + "\n")
				}
			}
		}

		writeErr = true
		for _, diff := range pDiffs {
			if diff.Error == ErrValuesDontMatch {
				if writeErr {
					prettyPrint.WriteString(ErrValuesDontMatch.Error() + "\n")
					writeErr = false
				}
				if diff.Left != "" {
					prettyPrint.WriteString("+" + diff.Left + "\n")
				}
				if diff.Right != "" {
					prettyPrint.WriteString("-" + diff.Right + "\n")
				}
			}
		}

		prettyPrint.WriteString("\n")
	}
	return prettyPrint.String()
}

func groupByPath(diffs []Diff) map[string][]Diff {
	var pathDiffs = make(map[string][]Diff)

	for i := range diffs {
		pathDiffs[diffs[i].Path] = append(pathDiffs[diffs[i].Path], diffs[i])
	}
	return pathDiffs
}
