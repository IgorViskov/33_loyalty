package appmath

import (
	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToNullable(t *testing.T) {
	nd := ToNullable(decimal.Decimal{})
	assert.NotNil(t, nd)
}

func TestParseDecimal(t *testing.T) {
	d, e := ParseDecimal("123.123")
	assert.Nil(t, e)
	assert.NotNil(t, d)
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		value    interface{}
		actual   int
		negative bool
		name     string
	}{
		{
			value:  "100",
			actual: 100,
			name:   "parse string positive",
		},
		{
			value:  200,
			actual: 200,
			name:   "parse int positive",
		},
		{
			value:    256.123,
			negative: true,
			name:     "parse overs negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			val, err := ParseInt(tt.value)
			if !tt.negative {
				assert.NoError(t, err)
				assert.Equal(t, tt.actual, val)
			} else {
				assert.Error(t, err)
			}

		})
	}
}
