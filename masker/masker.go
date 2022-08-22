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

// maskStruct analyzes a struct and applies masking if necessary
// it only masks fields that have been tagged as 'sensitive:"true"'
// e.g.
// type Example struct {
//		A string // no masking
//      B string `sensitive:"true"` // masking (partial)
//      C string `sensitive:"true,full"` // masking (full)
//}
func (m *masker) maskStruct(s interface{}) (interface{}, error) {
	if s == nil {
		return nil, errors.New("cannot mask nil value")
	}

	var (
		stReflect = reflect.ValueOf(s)
		st        = reflect.TypeOf(s)
		// create a new var with the same type of the struct
		cloned = reflect.New(st)
	)

	// iterate over the struct fields
	for i := 0; i < stReflect.NumField(); i++ {
		// get the field value
		value := stReflect.Field(i)
		// get the "sensitive" tag which indicates its value must be masked
		maskTag := stReflect.Type().Field(i).Tag.Get("sensitive")

		// apply different strategies depending on the field type
		// only apply masking if the mask tag is true
		switch value.Type().Kind() {
		case reflect.String:
			if maskTag != "" && maskTag != "false" {
				tmp := m.applyMask(interfaceToString(value), maskTag)
				// set the masked value to the corresponding field of the cloned struct
				cloned.Elem().Field(i).SetString(reflect.ValueOf(tmp).String())
				continue
			}
			// set the same value to the corresponding field of the cloned struct
			cloned.Elem().Field(i).Set(value)
			continue
		case reflect.Struct:
			// if field type is struct, then make a recursive call to maskStruct
			// and set the result to the corresponding field of the cloned struct
			tmp, err := m.maskStruct(value.Interface())
			if err != nil {
				return nil, err
			}
			cloned.Elem().Field(i).Set(reflect.ValueOf(tmp).Elem())
			continue
		case reflect.Slice:
			if value.IsNil() {
				continue
			}
			// []string -> apply masking to every item (if necessary)
			if value.Type().Elem().Kind() == reflect.String {
				values := value.Interface().([]string)
				nv := make([]string, len(values))
				for i, val := range values {
					if maskTag != "" && maskTag != "false" {
						nv[i] = m.applyMask(val, maskTag)
					} else {
						nv[i] = val
					}
				}
				cloned.Elem().Field(i).Set(reflect.ValueOf(nv))
				continue
			}
			// []struct {...} -> apply masking using recursive calls to maskStruct
			if value.Type().Elem().Kind() == reflect.Struct {
				nv := reflect.MakeSlice(value.Type(), 0, value.Len())
				for j, k := 0, value.Len(); j < k; j++ {
					tmp, err := m.maskStruct(value.Index(j).Interface())
					if err != nil {
						return nil, err
					}
					nv = reflect.Append(nv, reflect.ValueOf(tmp).Elem())
				}
				cloned.Elem().Field(i).Set(nv)
				continue
			}
		default:
			// for other types nothing should be masked
			cloned.Elem().Field(i).Set(value)
			continue
		}
	}

	// return
	return cloned.Interface(), nil
}

// applyMask apply masking based on the masking tag value
// e.g. sanitize:"true,full" -> apply full masking
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

// sanitize apply masking to a string using a provided mask character. Full-masking: replace the whole string with the
// mask character. Partial-masking: replace the 80% of the string with the masking character.
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
