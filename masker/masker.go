package masker

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	// FullMask indicates to apply masking
	FullMask = iota
	// PartialMask indicates to apply partial masking
	PartialMask
	defaultMask = "*"
)

// Masker provides several methods to mask sensitive data
type Masker interface {
	MaskSensitiveData(v interface{}) (interface{}, error)
	Sanitize(str string, level int) string
	SetMask(c string)
}

type masker struct {
	mask string
}

// New returns a new Masker instance
func New() Masker {
	return &masker{
		mask: defaultMask,
	}
}

// MaskSensitiveData auto-detects input type and masks sensitive data
func (m *masker) MaskSensitiveData(v interface{}) (interface{}, error) {
	d := reflect.ValueOf(v)

	switch d.Kind() {
	case reflect.Struct:
		return m.maskStruct(v)
	default:
		return v, nil
	}
}

func (m *masker) maskStruct(s interface{}) (interface{}, error) {
	if s == nil {
		return nil, errors.New("cannot mask nil value")
	}

	stReflect := reflect.ValueOf(s)
	var (
		st   = reflect.TypeOf(s)
		tmpr = reflect.New(st)
	)

	for i := 0; i < stReflect.NumField(); i++ {
		maskTag := stReflect.Type().Field(i).Tag.Get("sensitive")
		value := stReflect.Field(i)

		switch value.Type().Kind() {
		case reflect.String:
			if maskTag != "" && maskTag != "false" {
				tmp := m.applyMask(interfaceToString(value), maskTag)
				tmpr.Elem().Field(i).SetString(reflect.ValueOf(tmp).String())
				continue
			}
			tmpr.Elem().Field(i).Set(value)
			continue
		case reflect.Struct:
			tmp, err := m.maskStruct(value.Interface())
			if err != nil {
				return nil, err
			}
			tmpr.Elem().Field(i).Set(reflect.ValueOf(tmp).Elem())
			continue
		case reflect.Slice:
			if value.IsNil() {
				continue
			}
			if value.Type().Elem().Kind() == reflect.String {
				vals := value.Interface().([]string)
				nv := make([]string, len(vals))
				for i, val := range vals {
					if maskTag != "" && maskTag != "false" {
						nv[i] = m.applyMask(val, maskTag)
					} else {
						nv[i] = val
					}
				}
				tmpr.Elem().Field(i).Set(reflect.ValueOf(nv))
				continue
			}
			if value.Type().Elem().Kind() == reflect.Struct {
				nv := reflect.MakeSlice(value.Type(), 0, value.Len())
				for j, k := 0, value.Len(); j < k; j++ {
					tmp, err := m.maskStruct(value.Index(j).Interface())
					if err != nil {
						return nil, err
					}
					nv = reflect.Append(nv, reflect.ValueOf(tmp).Elem())
				}
				tmpr.Elem().Field(i).Set(nv)
				continue
			}
		default:
			tmpr.Elem().Field(i).Set(value)
			continue
		}
	}

	return tmpr.Interface(), nil
}

func (m *masker) applyMask(v, maskTag string) string {
	var (
		args      = strings.Split(maskTag, ",")
		maskLevel = parseMaskLevel(args)
	)
	return sanitize(v, m.mask, maskLevel)
}

func parseMaskLevel(tag []string) int {
	if len(tag) > 1 && tag[1] == "full" {
		return FullMask
	}
	return PartialMask
}

func sanitize(v, mask string, level int) string {
	length := len(v)
	if length == 0 {
		return v
	}

	if level == FullMask {
		return strings.Repeat(mask, length)
	}

	nonMaskedChars := ((20 * length) / 100) + 1
	return v[0:nonMaskedChars] + strings.Repeat(mask, length-nonMaskedChars)
}

// Sanitize apply masking with level to a string
func (m *masker) Sanitize(str string, level int) string {
	if level != FullMask && level != PartialMask {
		level = PartialMask
	}

	return sanitize(str, m.mask, level)
}

// SetMask sets the mask character
func (m *masker) SetMask(c string) {
	m.mask = c
}

func interfaceToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
