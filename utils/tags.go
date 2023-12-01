package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type TagsParser interface {
	TagName() string
	ParseTags(fieldVal reflect.Value, tags []string) error
}

func ParseFields(structPtr any, parser TagsParser) error {
	ptrVal := reflect.ValueOf(structPtr)
	if ptrVal.Kind() != reflect.Pointer || ptrVal.Elem().Kind() != reflect.Struct {
		return errors.New("value is not a pointer to struct")
	}
	v := ptrVal.Elem()
	t := reflect.TypeOf(structPtr).Elem()
	for i := 0; i < t.NumField(); i++ {
		tags := strings.Split(t.Field(i).Tag.Get(parser.TagName()), ",")
		for j := range tags {
			tags[j] = strings.TrimSpace(tags[j])
		}
		if err := parser.ParseTags(v.Field(i), tags); err != nil {
			return fmt.Errorf("error parsing tags: %w", err)
		}
	}
	return nil
}
