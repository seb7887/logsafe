package logger

import (
	"github.com/stretchr/testify/assert"
	"seb7887/logsafe/masker"
	"testing"
)

func TestFormatter(t *testing.T) {
	var (
		m       = masker.New()
		example = struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{
			Name:  "test",
			Value: "value",
		}
		specs = []struct {
			name string
			in   []interface{}
			out  string
		}{
			{
				name: "should format a single string value",
				in:   []interface{}{"test"},
				out:  "test",
			},
			{
				name: "should format multiple values",
				in:   []interface{}{"example: ", example},
				out:  "example: {\"name\":\"test\",\"value\":\"value\"}",
			},
		}
	)

	for _, spec := range specs {
		res := formatArgs(m, spec.in...)

		assert.Equal(t, spec.out, res)
	}
}
