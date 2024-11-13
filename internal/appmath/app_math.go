package appmath

import (
	"fmt"
	"github.com/govalues/decimal"
	"strconv"
)

func ToNullable(d decimal.Decimal) decimal.NullDecimal {
	return decimal.NullDecimal{
		Decimal: d,
		Valid:   true,
	}
}

func ParseDecimal(val interface{}) (decimal.Decimal, error) {
	switch v := val.(type) {
	case int:
		return decimal.NewFromInt64(int64(v), 0, 0)
	case string:
		return decimal.Parse(v)
	case float32:
		return decimal.NewFromFloat64(float64(v))
	case float64:
		return decimal.NewFromFloat64(v)
	default:
		return decimal.Decimal{}, fmt.Errorf("cannot parse decimal of type %T", v)
	}
}

func ParseInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot parse int of type %T", v)
	}
}
