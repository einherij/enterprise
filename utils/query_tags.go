package utils

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

const (
	tagName = "query"
)

type QueryUnmarshaler interface {
	UnmarshalQuery(param string) (err error)
}

type QueryMarshaller interface {
	MarshalQuery() (param string, err error)
}

func NewQueryDecoder(query url.Values) *QueryDecoder {
	return &QueryDecoder{
		query:      query,
		parseFuncs: make(map[string]func(string) (any, error)),
	}
}

type QueryDecoder struct {
	query      url.Values
	parseFuncs map[string]func(string) (any, error) // map[type]parseFunc
}

func (qp *QueryDecoder) DecodeQuery(dst any) error {
	if err := ParseFields(dst, qp); err != nil {
		return fmt.Errorf("error parsing query fields: %w", err)
	}
	return nil
}

func (qp *QueryDecoder) AddTypeParsing(typeSample any, parseFunc func(string) (any, error)) {
	qp.parseFuncs[reflect.ValueOf(typeSample).Type().String()] = parseFunc
}

func (qp *QueryDecoder) ParseTags(fieldVal reflect.Value, tags []string) error {
	strType := fieldVal.Type().String()
	for _, tag := range tags {
		if q := qp.query.Get(tag); q != "" {
			if fieldVal.CanAddr() && fieldVal.Addr().CanInterface() {
				if parser, ok := fieldVal.Addr().Interface().(QueryUnmarshaler); ok {
					if err := parser.UnmarshalQuery(q); err != nil {
						return fmt.Errorf("error parsing query using type method: %w", err)
					}
					return nil
				}
			}

			if pf, ok := qp.parseFuncs[strType]; ok {
				v, err := pf(q)
				if err != nil {
					return fmt.Errorf("error parsing %s value %q: %s", strType, q, err)
				}
				fieldVal.Set(reflect.ValueOf(v))
				return nil
			}

			switch fieldVal.Type().String() {
			case "string":
				fieldVal.Set(reflect.ValueOf(q))
			case "int":
				i, err := strconv.Atoi(q)
				if err != nil {
					return fmt.Errorf("error parsing int value %q: %s", q, err)
				}
				fieldVal.Set(reflect.ValueOf(i))
			case "float64":
				f, err := strconv.ParseFloat(q, 64)
				if err != nil {
					return fmt.Errorf("error parsing float value %q: %s", q, err)
				}
				fieldVal.Set(reflect.ValueOf(f))
			case "time.Time":
				i, err := strconv.ParseInt(q, 10, 64)
				if err != nil {
					return fmt.Errorf("error parsing timestamp value %q: %s", q, err)
				}
				t := time.Unix(i, 0)
				fieldVal.Set(reflect.ValueOf(t))
			default:
				return fmt.Errorf("no parsing function for type %s", fieldVal.Type())
			}

			return nil
		}
	}
	return nil
}

func (qp *QueryDecoder) TagName() string {
	return tagName
}

type QueryEncoder struct {
	query url.Values
}

func NewQueryEncoder(query url.Values) *QueryEncoder {
	return &QueryEncoder{
		query: query,
	}
}

func (qe *QueryEncoder) EncodeQuery(src any) error {
	if err := ParseFields(src, qe); err != nil {
		return fmt.Errorf("error encoding query fields: %w", err)
	}
	return nil
}

func (qe *QueryEncoder) ParseTags(fieldVal reflect.Value, tags []string) error {
	for _, tag := range tags {
		var (
			strVal string
		)
		if tag == "-" {
			return nil
		}
		if fieldVal.CanAddr() && fieldVal.Addr().CanInterface() {
			if parser, ok := fieldVal.Addr().Interface().(QueryMarshaller); ok {
				if p, err := parser.MarshalQuery(); err != nil {
					return fmt.Errorf("error parsing query using type method: %w", err)
				} else {
					strVal = p
				}
			}
		}
		if strVal == "" {
			strVal = fmt.Sprintf("%v", fieldVal)
		}
		if strVal != "" && strVal != "0" {
			qe.query.Set(tag, strVal)
		}

		break
	}
	return nil
}

func (qe *QueryEncoder) TagName() string {
	return tagName
}
