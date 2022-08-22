package masker

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type Address struct {
	Number int    `json:"number"`
	Street string `json:"street" sensitive:"true"`
}

type TestStruct struct {
	ID             string    `json:"id" sensitive:"true"`
	Username       string    `json:"full_name"`
	CardNumber     string    `json:"card_number" sensitive:"true"`
	Password       string    `json:"password" sensitive:"true,full"`
	Amount         float64   `json:"amount"`
	Address        Address   `json:"address"`
	IPs            []string  `json:"ips" sensitive:"true"`
	ChildAddresses []Address `json:"child_addresses"`
}

func TestMasker_MaskSensitiveData(t *testing.T) {
	var (
		m          = New()
		testStruct = TestStruct{
			ID:         "nuih",
			Username:   "test",
			CardNumber: "4312123453459909",
			Password:   "abc1234",
			Amount:     123.34,
			Address: Address{
				Number: 123,
				Street: "Something",
			},
			IPs: []string{"122345", "123454"},
			ChildAddresses: []Address{
				{
					Number: 222,
					Street: "child1",
				},
				{
					Number: 111,
					Street: "child2",
				},
			},
		}
		resStruct = &TestStruct{
			ID:         "n***",
			Username:   "test",
			CardNumber: "4312************",
			Password:   "*******",
			Amount:     123.34,
			Address: Address{
				Number: 123,
				Street: "So*******",
			},
			IPs: []string{"12****", "12****"},
			ChildAddresses: []Address{
				{
					Number: 222,
					Street: "ch****",
				},
				{
					Number: 111,
					Street: "ch****",
				},
			},
		}
		specs = []struct {
			name string
			in   interface{}
			out  interface{}
			err  error
		}{
			{
				name: "should mask sensitive data using filter values",
				in:   testStruct,
				out:  resStruct,
				err:  nil,
			},
		}
	)

	for _, spec := range specs {
		res, err := m.MaskSensitiveData(spec.in)

		if spec.err != nil {
			require.Error(t, err, "should return error")
		} else {
			require.NoError(t, err, "should not return error")
			assert.Equal(t, spec.out, res)
		}
	}
}

func TestMasker_SetMask(t *testing.T) {
	var (
		m  = New()
		in = "abc123"
	)

	out := m.Sanitize(in, PartialMask)
	assert.Equal(t, "ab****", out)

	m.SetMask("#")

	out = m.Sanitize(in, PartialMask)
	assert.Equal(t, "ab####", out)
}
